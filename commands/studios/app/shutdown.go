package app

import (
	"fmt"
	"os"

	"github.com/hofstadter-io/hof/pkg/studios/app"
	"github.com/spf13/cobra"
)

var ShutdownLong = `Shutdowns the App`

var ShutdownCmd = &cobra.Command{

	Use: "shutdown",

	Short: "Shutdowns the App",

	Long: ShutdownLong,

	Run: func(cmd *cobra.Command, args []string) {

		var name string
		if 0 < len(args) {
			name = args[0]
		}

		/*
			fmt.Println("hof app shutdown:",
				name,
			)
		*/

		err := app.Shutdown(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
