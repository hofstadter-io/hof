package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg"
)

var cmdLong = `run commands defined in _tool.cue files`

var CmdCmd = &cobra.Command{

	Use: "cmd [flags] [cmd] [args]",

	Short: "run commands defined in _tool.cue files",

	Long: cmdLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		flags := []string{}
		msg, err := pkg.Cmd(flags, args, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(msg)

	},
}
