package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type spinnerModel struct {
	spinner spinner.Model
	title   string
	done    bool
	err     error
}

type spinnerDoneMsg struct{ err error }

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinnerDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m spinnerModel) View() string {
	if m.done {
		return ""
	}
	return fmt.Sprintf("%s %s\n", m.spinner.View(), m.title)
}

// Spin runs fn in a background goroutine while showing a spinner with title.
// Writes spinner output to stderr.
func Spin(title string, fn func() error) error {
	if !IsTerminal() {
		fmt.Fprintln(os.Stderr, title)
		return fn()
	}

	s := spinner.New()
	s.Spinner = spinner.Meter
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := spinnerModel{spinner: s, title: title}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))

	doneCh := make(chan error, 1)
	go func() {
		doneCh <- fn()
	}()

	go func() {
		err := <-doneCh
		p.Send(spinnerDoneMsg{err: err})
	}()

	result, err := p.Run()
	if err != nil {
		return err
	}
	return result.(spinnerModel).err
}
