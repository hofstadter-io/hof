package app

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/hofstadter-io/hof/lib/app"
	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   delete
// Usage:  delete <name or id>
// Parent: app

var DeleteLong = `Delete an App and all associated data`

var DeleteCmd = &cobra.Command{

	Use: "delete <name or id>",

	Short: "Delete an App",

	Long: DeleteLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In deleteCmd", "args", args)
		// Argument Parsing
		// [0]name:   name
		//     help:
		//     req'd:  false

		var name string

		if 0 < len(args) {

			name = args[0]
		}

		fmt.Println("hof app delete:",
			name,
		)

		err := app.Delete(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	// add sub-commands to this command when present

}
