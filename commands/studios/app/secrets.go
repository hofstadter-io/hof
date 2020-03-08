package app

import (
	"fmt"

	// custom imports

	// infered imports

	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/app"
)

// Tool:   hof
// Name:   secrets
// Usage:  secrets
// Parent: app

var SecretsLong = `Set the App Secrets`

var SecretsCmd = &cobra.Command{

	Use: "secrets",

	Short: "Set the App Secrets",

	Long: SecretsLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In secretsCmd", "args", args)
		// Argument Parsing

		// fmt.Println("hof app secrets:")

		err := app.Secrets()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
