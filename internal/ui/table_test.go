package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// feed drives the model through a sequence of messages, returning the result.
func feed(m tableModel, msgs ...tea.Msg) tableModel {
	for _, msg := range msgs {
		model, _ := m.Update(msg)
		m = model.(tableModel)
	}
	return m
}

// typeStr simulates the user typing s one rune at a time.
func typeStr(m tableModel, s string) tableModel {
	for _, r := range s {
		m = feed(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	return m
}

// filterTestModel builds a select-mode model with a hidden IP column, so tests
// can verify filtering matches fields that are not displayed.
func filterTestModel() tableModel {
	columns := []string{"Name"} // only Name is shown; IPs are hidden
	rows := [][]string{
		{"web-prod"},
		{"db-prod"},
		{"web-staging"},
	}
	filterKeys := []string{
		"i-aaa web-prod 10.0.0.1 1.2.3.4",
		"i-bbb db-prod 10.0.0.2 ",
		"i-ccc web-staging 10.0.0.3 ",
	}
	return newTableModel(columns, rows, filterKeys, true)
}

func TestTableFilterNarrowsRows(t *testing.T) {
	m := filterTestModel()
	if got := len(m.table.Rows()); got != 3 {
		t.Fatalf("initial rows = %d, want 3", got)
	}

	m = typeStr(m, "web")
	if got := len(m.table.Rows()); got != 2 {
		t.Errorf("after typing 'web' rows = %d, want 2", got)
	}
	for _, r := range m.table.Rows() {
		if r[0] != "web-prod" && r[0] != "web-staging" {
			t.Errorf("unexpected filtered row %v", r)
		}
	}
}

func TestTableFilterMatchesHiddenColumn(t *testing.T) {
	m := filterTestModel()
	// 10.0.0.2 is only in the (hidden) filter key of db-prod.
	m = typeStr(m, "10.0.0.2")
	rows := m.table.Rows()
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	if rows[0][0] != "db-prod" {
		t.Errorf("filtered row = %v, want db-prod", rows[0])
	}
}

func TestTableClearFilterRestoresRows(t *testing.T) {
	m := filterTestModel()
	m = typeStr(m, "web")
	if len(m.table.Rows()) != 2 {
		t.Fatalf("precondition failed: rows = %d, want 2", len(m.table.Rows()))
	}
	// Backspace three times to clear "web".
	m = feed(m,
		tea.KeyMsg{Type: tea.KeyBackspace},
		tea.KeyMsg{Type: tea.KeyBackspace},
		tea.KeyMsg{Type: tea.KeyBackspace},
	)
	if got := len(m.table.Rows()); got != 3 {
		t.Errorf("after clearing filter rows = %d, want 3", got)
	}
}

func TestTableEnterSelectsFilteredRow(t *testing.T) {
	m := filterTestModel()
	m = typeStr(m, "10.0.0.2") // narrows to db-prod
	m = feed(m, tea.KeyMsg{Type: tea.KeyEnter})
	if !m.quitting {
		t.Fatalf("expected quitting after enter")
	}
	if len(m.selected) != 1 || m.selected[0] != "db-prod" {
		t.Errorf("selected = %v, want [db-prod]", m.selected)
	}
}

func TestTableEnterWithNoMatchDoesNotSelect(t *testing.T) {
	m := filterTestModel()
	m = typeStr(m, "zzz") // matches nothing
	if len(m.table.Rows()) != 0 {
		t.Fatalf("rows = %d, want 0", len(m.table.Rows()))
	}
	m = feed(m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.quitting {
		t.Errorf("should not quit when nothing matches")
	}
	if m.selected != nil {
		t.Errorf("selected = %v, want nil", m.selected)
	}
}

func TestTableEscCancels(t *testing.T) {
	m := filterTestModel()
	m = typeStr(m, "web")
	m = feed(m, tea.KeyMsg{Type: tea.KeyEsc})
	if !m.quitting {
		t.Errorf("expected quitting after esc")
	}
	if m.selected != nil {
		t.Errorf("selected = %v, want nil after cancel", m.selected)
	}
}

func TestTableQIsFilterCharWhenFilterable(t *testing.T) {
	m := filterTestModel()
	// 'q' should type into the filter, not quit.
	m = feed(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if m.quitting {
		t.Errorf("'q' should not quit while filtering")
	}
	if m.filter.Value() != "q" {
		t.Errorf("filter value = %q, want %q", m.filter.Value(), "q")
	}
}

func TestTableQQuitsWhenNotFilterable(t *testing.T) {
	// View-only table (no filter keys): 'q' quits.
	m := newTableModel([]string{"Name"}, [][]string{{"web-prod"}}, nil, false)
	if m.filterable {
		t.Fatalf("expected non-filterable model")
	}
	m = feed(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if !m.quitting {
		t.Errorf("expected 'q' to quit a non-filterable table")
	}
}
