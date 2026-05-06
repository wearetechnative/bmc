package ui

import (
	"fmt"
	"os"

	btable "github.com/charmbracelet/bubbles/table"
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
	selected   []string
	quitting   bool
	selectMode bool
	totalRows  int
	wantHeight int
}

func (m tableModel) Init() tea.Cmd { return nil }

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.selectMode {
				m.selected = m.table.SelectedRow()
				m.quitting = true
				return m, tea.Quit
			}
		case "ctrl+c", "esc", "q":
			m.quitting = true
			return m, tea.Quit
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
		hint = "↑↓ scroll  enter select  q cancel"
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
	return runTable(columns, rows, false, nil)
}

// SelectFromTable displays a table and returns the selected row values.
// Returns nil if the user cancelled.
func SelectFromTable(columns []string, rows [][]string) ([]string, error) {
	if !IsTerminal() {
		return selectPlainTable(columns, rows)
	}
	var selected []string
	err := runTable(columns, rows, true, &selected)
	return selected, err
}

func runTable(columns []string, rows [][]string, selectMode bool, out *[]string) error {
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

	// Desired height: all rows, capped at terminal height - 5 (border + header + footer + margin).
	wantHeight := len(rows) + 1
	height := wantHeight

	var p *tea.Program
	if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		if _, termH, err := term.GetSize(int(tty.Fd())); err == nil && termH > 0 {
			if max := termH - 5; max > 0 && height > max {
				height = max
			}
		}
		t := btable.New(
			btable.WithColumns(cols),
			btable.WithRows(tableRows),
			btable.WithFocused(true),
			btable.WithHeight(height),
		)
		s := btable.DefaultStyles()
		s.Header = s.Header.Bold(true)
		s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
		t.SetStyles(s)
		m := tableModel{table: t, selectMode: selectMode, totalRows: len(rows), wantHeight: wantHeight}
		p = tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))
	} else {
		t := btable.New(
			btable.WithColumns(cols),
			btable.WithRows(tableRows),
			btable.WithFocused(true),
			btable.WithHeight(height),
		)
		s := btable.DefaultStyles()
		s.Header = s.Header.Bold(true)
		s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
		t.SetStyles(s)
		m := tableModel{table: t, selectMode: selectMode, totalRows: len(rows), wantHeight: wantHeight}
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
