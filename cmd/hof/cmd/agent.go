package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	agentcmd "github.com/hofstadter-io/hof/lib/agent/cmd"
)

var agentLong = `Run an agent`

func init() {

	flags.SetupAgentFlags(AgentCmd.Flags(), &(flags.AgentFlags))

}

func AgentRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = agentcmd.Run(args, flags.RootPflags, flags.AgentFlags)

	return err
}

var AgentCmd = &cobra.Command{

	Use: "agent [args]",

	Short: "run an agent",

	Long: agentLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = AgentRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := AgentCmd.HelpFunc()
	ousage := AgentCmd.UsageFunc()

	help := func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath() + " help")

		if extra(cmd) {
			return
		}
		ohelp(cmd, args)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		return ousage(cmd)
	}

	thelp := func(cmd *cobra.Command, args []string) {
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		return usage(cmd)
	}
	AgentCmd.SetHelpFunc(thelp)
	AgentCmd.SetUsageFunc(tusage)

}
