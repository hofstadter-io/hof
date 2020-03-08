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
// Name:   push
// Usage:  push
// Parent: function

var PushLong = `Uploads the local copy and makes it the latest copy in Studios`

var PushCmd = &cobra.Command{

	Use: "push <function path>",

	Short: "Send and make the latest version on Studios",

	Long: PushLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In pushCmd", "args", args)

		// Argument Parsing

		fmt.Println("hof function push: ")

		err := fns.Push()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
