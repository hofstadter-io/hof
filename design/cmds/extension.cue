package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

ExtensionCommand: schema.Command & {
	Name:  "extension"
	Usage: "extension [args]"
	Short: "run the extension server"
	Long:  "run the extension server"
}
