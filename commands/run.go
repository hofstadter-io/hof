package commands

import (

	// hello... something might need to go here

	"github.com/spf13/cobra"

	"fmt"

	"os"

	"github.com/hofstadter-io/hof/lib"
)

var runLong = `run commands defined by HofCmd. Falls back on cue commands, which also falls back to their own run system`

var RunCmd = &cobra.Command{

	Use: "run [flags] [cmd] [args]",

	Short: "run commands defined by HofCmd",

	Long: runLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		flags := []string{}
		msg, err := lib.Cmd(flags, args, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(msg)

	},
}
