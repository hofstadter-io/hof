package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/app"
)

var DeleteLong = `Delete an App and all associated data`

var DeleteCmd = &cobra.Command{

	Use: "delete <name or id>",

	Short: "Delete an App",

	Long: DeleteLong,

	Run: func(cmd *cobra.Command, args []string) {

		var name string
		if 0 < len(args) {
			name = args[0]
		}

		/*
			fmt.Println("hof app delete:",
				name,
			)
		*/

		err := app.Delete(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
