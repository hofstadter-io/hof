package database

import (
	"fmt"
	"os"

	"github.com/hofstadter-io/hof/pkg/studios/database"
	"github.com/spf13/cobra"
)

var StatusLong = `Get the status of your DB`

var StatusCmd = &cobra.Command{

	Use: "status",

	Short: "Get the status of your DB",

	Long: StatusLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof db status:")

		err := database.Status()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
