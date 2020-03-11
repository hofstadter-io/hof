package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/app"
)

var UpdateLong = `Updates the Application runtime when a new version is available. If you omit version, the most recent, stable version will be applied.`

var UpdateCmd = &cobra.Command{

	Use: "update [version]",

	Short: "Updates the Application runtime",

	Long: UpdateLong,

	Run: func(cmd *cobra.Command, args []string) {

		var version string

		if 0 < len(args) {
			version = args[0]
		}

		/*
			fmt.Println("hof app update:",
				version,
			)
		*/

		err := app.Update(version)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
