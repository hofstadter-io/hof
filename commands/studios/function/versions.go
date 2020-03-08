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
// Name:   versions
// Usage:  versions
// Parent: function

var VersionsLong = `Get the supported runtime versions for Hofstadter Studios`

var VersionsCmd = &cobra.Command{

	Use: "versions",

	Short: "Get the runtime versions",

	Long: VersionsLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In versionsCmd", "args", args)
		// Argument Parsing

		/*
			fmt.Println("hof function versions:")
		*/

		err := fns.Versions()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
