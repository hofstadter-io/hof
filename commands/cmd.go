package commands

import (

	// hello... something might need to go here

	"github.com/spf13/cobra"

	"fmt"

	"os"

	"github.com/hofstadter-io/hof/lib"
)

var cmdLong = `run commands defined in _tool.cue files`

var CmdCmd = &cobra.Command{

	Use: "cmd [flags] [cmd] [args]",

	Short: "run commands defined in _tool.cue files",

	Long: cmdLong,

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
