package ui

import (
	"fmt"
	"os"
	"strings"

	btable "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lgtable "github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

var tableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableModel struct {
	table      btable.Model
	filter     textinput.Model
	allRows    []btable.Row // full, unfiltered row set
	filterKeys []string     // lowercased searchable string per allRows index
	filterable bool         // type-to-filter enabled (select mode with filter keys)
	selected   []string
	quitting   bool
	selectMode bool
	totalRows  int
	wantHeight int
}

func (m tableModel) Init() tea.Cmd { return nil }

// applyFilter recomputes the visible rows from the current filter query.
// An empty query shows all rows. SetRows clamps the cursor automatically.
func (m *tableModel) applyFilter() {
	q := strings.ToLower(strings.TrimSpace(m.filter.Value()))
	if q == "" {
		m.table.SetRows(m.allRows)
		return
	}
	var filtered []btable.Row
	for i, key := range m.filterKeys {
		if strings.Contains(key, q) {
			filtered = append(filtered, m.allRows[i])
		}
	}
	m.table.SetRows(filtered)
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.selectMode {
				// Nothing to select when the filter hides every row.
				if row := m.table.SelectedRow(); row != nil {
					m.selected = row
					m.quitting = true
					return m, tea.Quit
				}
				return m, nil
			}
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "up", "down", "pgup", "pgdown", "home", "end":
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		case "q":
			// When filtering, 'q' is a normal character; only quit otherwise.
			if !m.filterable {
				m.quitting = true
				return m, tea.Quit
			}
			fallthrough
		default:
			if m.filterable {
				var cmd tea.Cmd
				m.filter, cmd = m.filter.Update(msg)
				m.applyFilter()
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		if msg.Height > 0 {
			h := m.wantHeight
			if max := msg.Height - 5; max > 0 && h > max {
				h = max
			}
			m.table.SetHeight(h)
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	if m.quitting {
		return ""
	}
	hint := "↑↓ scroll  q quit"
	if m.selectMode {
		hint = "↑↓ scroll  enter select  esc cancel"
	}

	var header string
	if m.filterable {
		shown := len(m.table.Rows())
		footer := fmt.Sprintf("  %d/%d rows  type to filter  %s", shown, m.totalRows, hint)
		header = m.filter.View() + "\n"
		return header + tableStyle.Render(m.table.View()) + "\n" + footer + "\n"
	}

	footer := fmt.Sprintf("  %d rows  %s", m.totalRows, hint)
	return tableStyle.Render(m.table.View()) + "\n" + footer + "\n"
}

// ShowTable displays a table without selection and waits for the user to quit.
func ShowTable(columns []string, rows [][]string) error {
	if !IsTerminal() {
		printPlainTable(columns, rows)
		return nil
	}
	return runTable(columns, rows, nil, false, nil)
}

// SelectFromTable displays a table and returns the selected row values.
// Returns nil if the user cancelled. When filterKeys is non-nil (one lowercased
// searchable string per row), the picker supports real-time type-to-filter.
func SelectFromTable(columns []string, rows [][]string, filterKeys []string) ([]string, error) {
	if !IsTerminal() {
		return selectPlainTable(columns, rows)
	}
	var selected []string
	err := runTable(columns, rows, filterKeys, true, &selected)
	return selected, err
}

// newTableModel builds the interactive table model. Filtering is enabled when
// selectMode is true and filterKeys is non-nil (one lowercased searchable
// string per row). The table height defaults to fit all rows; callers may clamp
// it afterwards with m.table.SetHeight.
func newTableModel(columns []string, rows [][]string, filterKeys []string, selectMode bool) tableModel {
	filterable := selectMode && filterKeys != nil

	widths := make([]int, len(columns))
	for i, c := range columns {
		widths[i] = len(c)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	cols := make([]btable.Column, len(columns))
	for i, c := range columns {
		cols[i] = btable.Column{Title: c, Width: widths[i] + 2}
	}

	tableRows := make([]btable.Row, len(rows))
	for i, r := range rows {
		tableRows[i] = btable.Row(r)
	}

	ti := textinput.New()
	ti.Prompt = "filter: "
	ti.Placeholder = "type to filter"
	if filterable {
		ti.Focus()
	}

	// Desired height: all rows, capped at terminal height - 5 (border + header + footer + margin).
	wantHeight := len(rows) + 1

	t := btable.New(
		btable.WithColumns(cols),
		btable.WithRows(tableRows),
		btable.WithFocused(true),
		btable.WithHeight(wantHeight),
	)
	s := btable.DefaultStyles()
	s.Header = s.Header.Bold(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	t.SetStyles(s)

	return tableModel{
		table:      t,
		filter:     ti,
		allRows:    tableRows,
		filterKeys: filterKeys,
		filterable: filterable,
		selectMode: selectMode,
		totalRows:  len(rows),
		wantHeight: wantHeight,
	}
}

func runTable(columns []string, rows [][]string, filterKeys []string, selectMode bool, out *[]string) error {
	m := newTableModel(columns, rows, filterKeys, selectMode)

	var p *tea.Program
	if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(tty))
		if _, termH, err := term.GetSize(int(tty.Fd())); err == nil && termH > 0 {
			if max := termH - 5; max > 0 && m.wantHeight > max {
				m.table.SetHeight(max)
			}
		}
		p = tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))
	} else {
		p = tea.NewProgram(m, tea.WithOutput(os.Stderr))
	}

	result, err := p.Run()
	if err != nil {
		return err
	}
	if out != nil {
		*out = result.(tableModel).selected
	}
	return nil
}

// PrintTable writes a bordered, formatted table to stdout using lipgloss.
func PrintTable(columns []string, rows [][]string) {
	t := lgtable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		Headers(columns...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == lgtable.HeaderRow {
				return lipgloss.NewStyle().Bold(true)
			}
			return lipgloss.NewStyle()
		})
	for _, row := range rows {
		t = t.Row(row...)
	}
	fmt.Fprintln(os.Stdout, t.String())
}

func printPlainTable(columns []string, rows [][]string) {
	PrintTable(columns, rows)
}

func selectPlainTable(columns []string, rows [][]string) ([]string, error) {
	printPlainTable(columns, rows)
	var n int
	fmt.Fprint(os.Stderr, "Select row number: ")
	if _, err := fmt.Scan(&n); err != nil || n < 1 || n > len(rows) {
		return nil, nil
	}
	return rows[n-1], nil
}
