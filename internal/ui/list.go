package ui

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// ErrBack is returned by Choose when the user pressed ESC to go back.
// Callers that support multi-level navigation should check errors.Is(err, ErrBack)
// and re-present the previous menu. Ctrl+C returns ("", nil) instead.
var ErrBack = errors.New("ui: user navigated back")

// Item is a selectable list item.
type Item struct {
	Title string
	Desc  string
}

func (i Item) FilterValue() string { return i.Title }

// itemDelegate renders Item using the Title field directly,
// since Item has a field named Title which prevents adding a Title() method.
type itemDelegate struct{}

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                           { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(Item)
	if !ok {
		return
	}

	line := i.Title
	if i.Desc != "" {
		line += dimStyle.Render("  "+i.Desc)
	}

	if index == m.Index() {
		fmt.Fprint(w, selectedStyle.Render("> "+i.Title)+func() string {
			if i.Desc != "" {
				return dimStyle.Render("  " + i.Desc)
			}
			return ""
		}())
	} else {
		fmt.Fprint(w, normalStyle.Render("  "+line))
	}
}

type listModel struct {
	list       list.Model
	selected   string
	quitting   bool
	wentBack   bool // true when user pressed ESC (go back), false for Ctrl+C (cancel)
	wantHeight int  // desired height before terminal-size clamp
}

func (m listModel) Init() tea.Cmd { return nil }

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// When filtering, let the list handle Enter to commit the filter first
			if m.list.FilterState() == list.Filtering {
				break
			}
			if i, ok := m.list.SelectedItem().(Item); ok {
				m.selected = i.Title
			}
			m.quitting = true
			return m, tea.Quit
		case "esc":
			m.quitting = true
			m.wentBack = true
			return m, tea.Quit
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		if msg.Height > 0 {
			h := m.wantHeight
			if max := msg.Height - 2; h > max {
				h = max
			}
			m.list.SetHeight(h)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	if m.quitting {
		return ""
	}
	return m.list.View()
}

// Choose presents a filterable list and returns the selected item title.
// Returns ("", nil) if the user cancelled.
func Choose(header string, items []Item) (string, error) {
	if !IsTerminal() {
		return choosePlain(header, items)
	}

	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	var wantHeight int
	if len(items) <= 4 {
		wantHeight = len(items) + 3
	} else {
		wantHeight = len(items) + 6
		if wantHeight > 40 {
			wantHeight = 40
		}
	}

	// Clamp to terminal height and open /dev/tty for bubbletea I/O.
	height := wantHeight
	var p *tea.Program
	if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		if _, termH, err := term.GetSize(int(tty.Fd())); err == nil && termH > 0 {
			if max := termH - 2; height > max {
				height = max
			}
		}
		l := list.New(listItems, itemDelegate{}, 80, height)
		l.Title = header
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(true)
		l.SetShowHelp(len(items) > 4)
		l.SetShowPagination(len(items) > 4)
		m := listModel{list: l, wantHeight: wantHeight}
		p = tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))
	} else {
		l := list.New(listItems, itemDelegate{}, 80, height)
		l.Title = header
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(true)
		l.SetShowHelp(len(items) > 4)
		l.SetShowPagination(len(items) > 4)
		m := listModel{list: l, wantHeight: wantHeight}
		p = tea.NewProgram(m, tea.WithOutput(os.Stderr))
	}
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	m := result.(listModel)
	if m.wentBack {
		return "", ErrBack
	}
	return m.selected, nil
}

// choosePlain is a fallback for non-interactive environments.
func choosePlain(header string, items []Item) (string, error) {
	fmt.Fprintln(os.Stderr, header)
	for i, item := range items {
		fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, item.Title)
	}
	var n int
	fmt.Fprint(os.Stderr, "Select number: ")
	if _, err := fmt.Scan(&n); err != nil {
		return "", nil
	}
	if n < 1 || n > len(items) {
		return "", nil
	}
	return items[n-1].Title, nil
}

// IsTerminal returns true if stdin is a terminal.
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var _ = strings.Join // keep import
