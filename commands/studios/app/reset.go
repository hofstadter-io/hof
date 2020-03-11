package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/app"
)

var ResetLong = `Resets the App, because sometimes things get weird...`

var ResetCmd = &cobra.Command{

	Use: "reset",

	Short: "Reset the App",

	Long: ResetLong,

	Run: func(cmd *cobra.Command, args []string) {

		var name string
		if 0 < len(args) {
			name = args[0]
		}

		/*
			fmt.Println("hof app reset:",
				name,
			)
		*/

		err := app.Reset(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
