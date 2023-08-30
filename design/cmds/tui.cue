package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

TuiCommand: schema.Command & {
	Name:   "tui"
	Usage:  "tui"
	Short:  "hidden command for tui experiments"
	Long:   Short
	Hidden: true
}
