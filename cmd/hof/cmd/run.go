package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var runLong = `run commands defined by HofCmd. Falls back on cue commands, which also falls back to their own run system`

func RunRun(args []string) (err error) {

	flags, args := []string{}, []string{}
	msg, err := lib.Cmd(flags, args, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(msg)

	return err
}

var RunCmd = &cobra.Command{

	Use: "run [flags] [cmd] [args]",

	Short: "run commands defined by HofCmd",

	Long: runLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = RunRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
