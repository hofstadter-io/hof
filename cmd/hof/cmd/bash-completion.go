package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var BashCompletionLong = `Generate Bash completions

To load completion run

. <(hof bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(hof bash)
`

var BashCompletionCmd = &cobra.Command{

	Use: "bash-completion",

	Aliases: []string{
		"bash",
		"completions",
	},

	Short: "Generate  Bash completions",

	Long: BashCompletionLong,

	Run: func(cmd *cobra.Command, args []string) {

		RootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	RootCmd.AddCommand(BashCompletionCmd)
}
