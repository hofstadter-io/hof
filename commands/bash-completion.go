package commands

import (
	// "fmt"

	// custom imports

	// infered imports

	"os"
	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   bash-completion
// Usage:  bash-completion
// Parent: hof

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
		logger.Debug("In bash-completionCmd", "args", args)
		// Argument Parsing

		// fmt.Println("hof bash-completion:")

		RootCmd.GenBashCompletion(os.Stdout);
	},
}

func init() {
	RootCmd.AddCommand(BashCompletionCmd)
}

func init() {
	// add sub-commands to this command when present

}
