package container

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/hofstadter-io/hof/lib/crun"
	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   logs
// Usage:  logs [name]
// Parent: container

var LogsLong = `List the logs of your container.
If name is not specified, the current directory is used.
`

var LogsCmd = &cobra.Command{

	Use: "logs [name]",

	Short: "List the logs of your container",

	Long: LogsLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In logsCmd", "args", args)
		// Argument Parsing
		// [0]name:   name
		//     help:
		//     req'd:

		var name string

		if 0 < len(args) {

			name = args[0]
		}

		/*
			fmt.Println("hof containers logs:",
				name,
			)
		*/

		err := crun.Logs(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	// add sub-commands to this command when present

}
