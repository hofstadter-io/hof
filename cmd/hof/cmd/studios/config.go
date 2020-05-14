package cmdstudios

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios/config"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var configLong = `Work with Hofstadter Studios configs`

func ConfigRun(args []string) (err error) {

	return err
}

var ConfigCmd = &cobra.Command{

	Use: "config",

	Aliases: []string{
		"cfg",
	},

	Short: "Work with Hofstadter Studios configs",

	Long: configLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ConfigRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	hf := ConfigCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		hf(cmd, args)
	}
	ConfigCmd.SetHelpFunc(f)
	ConfigCmd.AddCommand(cmdconfig.ListCmd)
	ConfigCmd.AddCommand(cmdconfig.GetCmd)
	ConfigCmd.AddCommand(cmdconfig.CreateCmd)
	ConfigCmd.AddCommand(cmdconfig.UpdateCmd)
	ConfigCmd.AddCommand(cmdconfig.DeleteCmd)
}
