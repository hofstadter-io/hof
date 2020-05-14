package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var studiosLong = `  Hofstadter Studios makes it easy to develop and launch both
  hof-lang modules as well as pretty much any code or application`

var StudiosCmd = &cobra.Command{

	Use: "studios",

	Aliases: []string{
		"s",
	},

	Short: "commands for working with Hofstadter Studios",

	Long: studiosLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},
}

func init() {
	hf := StudiosCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		hf(cmd, args)
	}
	StudiosCmd.SetHelpFunc(f)
	StudiosCmd.AddCommand(cmdstudios.AppCmd)
	StudiosCmd.AddCommand(cmdstudios.DatabaseCmd)
	StudiosCmd.AddCommand(cmdstudios.ContainerCmd)
	StudiosCmd.AddCommand(cmdstudios.FunctionCmd)
	StudiosCmd.AddCommand(cmdstudios.ConfigCmd)
	StudiosCmd.AddCommand(cmdstudios.SecretCmd)
}
