package cmd

import (
	"strings"

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

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "", 0)

	},
}

func init() {

	help := TopicCmd.HelpFunc()
	usage := TopicCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "", 0)
		return usage(cmd)
	}
	TopicCmd.SetHelpFunc(thelp)
	TopicCmd.SetUsageFunc(tusage)

	TopicCmd.AddCommand(cmdtopic.WorkspaceCmd)

}
