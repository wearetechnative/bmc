package ui

import (
	"github.com/charmbracelet/huh"
)

// confirmHuh uses huh.Confirm for TTY environments.
func confirmHuh(prompt string) (bool, error) {
	var result bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(prompt).
				Value(&result),
		),
	).Run()
	if err == huh.ErrUserAborted {
		return false, nil
	}
	return result, err
}
