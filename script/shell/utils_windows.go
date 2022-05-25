//go:build windows
// +build windows

package shell

import (
	"github.com/chzyer/readline"
)

func clearScreen(s *Shell) error {
	return readline.ClearScreen(s.writer)
}
