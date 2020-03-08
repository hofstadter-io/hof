package secret

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/hofstadter-io/hof/lib/secret"
	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   list
// Usage:  list
// Parent: secret

var ListLong = `List your Studios secrets`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your secrets",

	Long: ListLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In listCmd", "args", args)
		// Argument Parsing

		/*
			fmt.Println("hof secret list:")
		*/

		err := secret.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
