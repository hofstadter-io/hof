package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/app"
)

var AvailableVersionLong = `Get the runtime version of your App`

var AvailableVersionCmd = &cobra.Command{

	Use: "available-versions",

	Aliases: []string{
		"versions",
		"avail-versions",
		"avail-vers",
		"vers",
	},

	Short: "Get the runtime version of your App",

	Long: AvailableVersionLong,

	Run: func(cmd *cobra.Command, args []string) {
		err := app.Versions()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

