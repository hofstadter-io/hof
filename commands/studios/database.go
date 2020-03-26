package studios

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/database"
)

var databaseLong = `Work with Hofstadter Studios databases`

var DatabaseCmd = &cobra.Command{

	Use: "database",

	Short: "Work with Hofstadter Studios databases",

	Long: databaseLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios database")

	},
}

func init() {
	DatabaseCmd.AddCommand(database.ListCmd)
	DatabaseCmd.AddCommand(database.GetCmd)
	DatabaseCmd.AddCommand(database.CreateCmd)
	DatabaseCmd.AddCommand(database.UpdateCmd)
	DatabaseCmd.AddCommand(database.StatusCmd)
	DatabaseCmd.AddCommand(database.SaveCmd)
	DatabaseCmd.AddCommand(database.RestoreCmd)
	DatabaseCmd.AddCommand(database.DeleteCmd)
}
