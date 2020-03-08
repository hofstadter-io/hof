package app

import (
	"fmt"
	"os"

	"github.com/hofstadter-io/hof/pkg/studios/app"
	"github.com/spf13/cobra"
)

var VersionLong = `Get the runtime version of your App`

var VersionCmd = &cobra.Command{

	Use: "version",

	Aliases: []string{
		"vers",
	},

	Short: "Get the runtime version of your App",

	Long: VersionLong,

	Run: func(cmd *cobra.Command, args []string) {

		var name string

		if 0 < len(args) {

			name = args[0]
		}

		/*
			fmt.Println("hof app version:",
				name,
			)
		*/

		err := app.Version(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
