package cmdstudios

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios/database"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var databaseLong = `Work with Hofstadter Studios databases`

func DatabaseRun(args []string) (err error) {

	return err
}

var DatabaseCmd = &cobra.Command{

	Use: "database",

	Short: "Work with Hofstadter Studios databases",

	Long: databaseLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},

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
	hf := DatabaseCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		as := strings.Join(args, "/")
		ga.SendGaEvent(c+"/help", as, 0)
		hf(cmd, args)
	}
	DatabaseCmd.SetHelpFunc(f)
	DatabaseCmd.AddCommand(cmddatabase.ListCmd)
	DatabaseCmd.AddCommand(cmddatabase.GetCmd)
	DatabaseCmd.AddCommand(cmddatabase.CreateCmd)
	DatabaseCmd.AddCommand(cmddatabase.UpdateCmd)
	DatabaseCmd.AddCommand(cmddatabase.StatusCmd)
	DatabaseCmd.AddCommand(cmddatabase.SaveCmd)
	DatabaseCmd.AddCommand(cmddatabase.RestoreCmd)
	DatabaseCmd.AddCommand(cmddatabase.DeleteCmd)
}
