package app

import (
	"fmt"

	// custom imports

	// infered imports

	"os"

	"github.com/hofstadter-io/hof/lib/app"
	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   available-version
// Usage:  available-versions
// Parent: app

var AvailableVersionLong = `Get the runtime version of your App`

var AvailableVersionCmd = &cobra.Command{

	Use: "available-versions",

	Aliases: []string{
		"versions",
		"avail-versions",
		"avail-vers",
		"vers",
	},

	Short: "Get the runtime version of your App",

	Long: AvailableVersionLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In available-versionCmd", "args", args)
		// Argument Parsing

		// fmt.Println("hof app available-version:")

		err := app.Versions()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
