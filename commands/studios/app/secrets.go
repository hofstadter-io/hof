package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/app"
)

var SecretsLong = `Set the App Secrets`

var SecretsCmd = &cobra.Command{

	Use: "secrets",

	Short: "Set the App Secrets",

	Long: SecretsLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof app secrets:")

		err := app.Secrets()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
