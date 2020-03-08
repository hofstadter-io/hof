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
// Name:   list
// Usage:  list
// Parent: container

var ListLong = `List your containers`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your containers",

	Long: ListLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In listCmd", "args", args)
		// Argument Parsing

		/*
			fmt.Println("hof containers list:")
		*/

		err := crun.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	// add sub-commands to this command when present

}
