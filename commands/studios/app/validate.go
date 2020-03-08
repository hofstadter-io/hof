package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/app"
)

var ValidateLong = `Validate your application or components of it`

var ValidateCmd = &cobra.Command{

	Use: "validate",

	Aliases: []string{
		"valid",
		"v",
	},

	Short: "Validate your application",

	Long: ValidateLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof app validate:")

		err := app.Validate()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
