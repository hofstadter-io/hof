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
// Name:   list
// Usage:  list
// Parent: function

var ListLong = `List your functions`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your functions",

	Long: ListLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In listCmd", "args", args)
		// Argument Parsing

		// fmt.Println("hof function list:")

		err := fns.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
