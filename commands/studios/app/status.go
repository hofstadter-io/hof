package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/app"
)

var StatusLong = `Get the status of your App`

var StatusCmd = &cobra.Command{

	Use: "status",

	Short: "Get the status of your App",

	Long: StatusLong,

	Run: func(cmd *cobra.Command, args []string) {

		var name string

		if 0 < len(args) {

			name = args[0]
		}

		/*
			fmt.Println("hof app status:",
				name,
			)
		*/

		err := app.Status(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
