package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/app"
)

var ListLong = `List app of your Apps`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List app of your Apps",

	Long: ListLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof app list:")

		err := app.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
