package cmdstudios

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios/database"
)

var databaseLong = `Work with Hofstadter Studios databases`

func DatabaseRun(args []string) (err error) {

	return err
}

var DatabaseCmd = &cobra.Command{

	Use: "database",

	Short: "Work with Hofstadter Studios databases",

	Long: databaseLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = DatabaseRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	DatabaseCmd.AddCommand(cmddatabase.ListCmd)
	DatabaseCmd.AddCommand(cmddatabase.GetCmd)
	DatabaseCmd.AddCommand(cmddatabase.CreateCmd)
	DatabaseCmd.AddCommand(cmddatabase.UpdateCmd)
	DatabaseCmd.AddCommand(cmddatabase.StatusCmd)
	DatabaseCmd.AddCommand(cmddatabase.SaveCmd)
	DatabaseCmd.AddCommand(cmddatabase.RestoreCmd)
	DatabaseCmd.AddCommand(cmddatabase.DeleteCmd)
}
