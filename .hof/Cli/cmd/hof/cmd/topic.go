package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/topic"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var topicLong = `Help for various topics and concepts`

func TopicRun(args []string) (err error) {

	return err
}

var TopicCmd = &cobra.Command{

	Use: "topic",

	Aliases: []string{
		"topics",
	},

	Short: "Help for various topics and concepts",

	Long: topicLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = TopicRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := TopicCmd.HelpFunc()
	usage := TopicCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		return usage(cmd)
	}
	TopicCmd.SetHelpFunc(thelp)
	TopicCmd.SetUsageFunc(tusage)

	TopicCmd.AddCommand(cmdtopic.WorkspaceCmd)

}
