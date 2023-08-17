package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

StCommand: schema.Command & {
	Name:  "st"
	Usage: "st"
	Short: "apply CUE transformations in bulk"
	Long:  Short

	Flags: [{
		Name:    "issue"
		Long:    "issue"
		Short:   "i"
		Type:    "bool"
		Default: "false"
		Help:    "create an issue (discussion is default)"
	}]
}
