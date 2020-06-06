package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/topic"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var topicLong = `help for various topics and concepts`

var TopicCmd = &cobra.Command{

	Use: "topic",

	Aliases: []string{
		"topics",
	},

	Short: "help for various topics and concepts",

	Long: topicLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {

	help := TopicCmd.HelpFunc()
	usage := TopicCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	TopicCmd.SetHelpFunc(thelp)
	TopicCmd.SetUsageFunc(tusage)

	TopicCmd.AddCommand(cmdtopic.WorkspaceCmd)

}
