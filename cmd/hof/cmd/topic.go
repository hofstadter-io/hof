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
	hf := TopicCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		hf(cmd, args)
	}
	TopicCmd.SetHelpFunc(f)
	TopicCmd.AddCommand(cmdtopic.WorkspaceCmd)
}
