package function

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/hofstadter-io/hof/lib/fns"

	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   logs
// Usage:  logs
// Parent: function

var LogsLong = `List the logs of your function`

var LogsCmd = &cobra.Command{

	Use: "logs",

	Short: "List the logs of your function",

	Long: LogsLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In logsCmd", "args", args)
		// Argument Parsing

		// fmt.Println("hof function logs:")

		err := fns.Logs()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
