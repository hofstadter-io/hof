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
// Name:   update
// Usage:  update [version]
// Parent: app

var UpdateLong = `Updates the Application runtime when a new version is available. If you omit version, the most recent, stable version will be applied.`

var UpdateCmd = &cobra.Command{

	Use: "update [version]",

	Short: "Updates the Application runtime",

	Long: UpdateLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In updateCmd", "args", args)
		// Argument Parsing
		// [0]name:   version
		//     help:
		//     req'd:

		var version string

		if 0 < len(args) {

			version = args[0]
		}

		/*
			fmt.Println("hof app update:",
				version,
			)
		*/

		err := app.Update(version)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
