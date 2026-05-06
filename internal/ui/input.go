package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	input    textinput.Model
	value    string
	quitting bool
}

func (m inputModel) Init() tea.Cmd { return textinput.Blink }

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.value = m.input.Value()
			m.quitting = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	if m.quitting {
		return ""
	}
	return m.input.View() + "\n"
}

// Input shows a text input prompt and returns the entered value.
// Returns ("", nil) if cancelled.
func Input(prompt string, echoMode bool) (string, error) {
	if !IsTerminal() {
		fmt.Fprint(os.Stderr, prompt)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			return scanner.Text(), nil
		}
		return "", nil
	}

	ti := textinput.New()
	ti.Placeholder = ""
	ti.Prompt = prompt + " "
	ti.Focus()
	if !echoMode {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	m := inputModel{input: ti}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	return result.(inputModel).value, nil
}

// Confirm asks a yes/no question. Returns true for yes.
func Confirm(prompt string) (bool, error) {
	if !IsTerminal() {
		fmt.Fprintf(os.Stderr, "%s [y/N]: ", prompt)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			ans := strings.ToLower(strings.TrimSpace(scanner.Text()))
			return ans == "y" || ans == "yes", nil
		}
		return false, nil
	}

	// Use huh for confirmation in TTY
	ans, err := confirmHuh(prompt)
	return ans, err
}
