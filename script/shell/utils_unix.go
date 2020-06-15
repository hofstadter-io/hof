// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd solaris

package shell

import (
	"github.com/chzyer/readline"
)

func clearScreen(s *Shell) error {
	_, err := readline.ClearScreen(s.writer)
	return err
}
