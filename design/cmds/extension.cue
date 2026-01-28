package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

ExtensionCommand: schema.Command & {
	Name:  "extension"
	Usage: "extension [args]"
	Short: "run the extension server"
	Long:  "run the extension server"

	Imports: [
		{Path: "github.com/hofstadter-io/hof/cmd/hof/flags"},
		{Path: "github.com/hofstadter-io/hof/lib/extension"},
		{Path: "github.com/hofstadter-io/hof/lib/runtime"},
	]

	PersistentPrerun:     true
	PersistentPrerunBody: "err = runtime.EnsureInfra()"

	// this runs the commands or naked `veg env` command
	Body: "err = extension.Run(args, flags.RootPflags)"
}
