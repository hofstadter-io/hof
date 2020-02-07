package commands

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg"
)

// Tool:   hof
// Name:   gen
// Usage:  gen <entrypoint>
// Parent: hof

var GenLong = `Generate a project starting from the file.`

var GenCmd = &cobra.Command{

	Use: "gen <entrypoint>",

	Aliases: []string{
		"g",
	},

	Short: "Generate a project.",

	Long: GenLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In genCmd", "args", args)
		// Argument Parsing

		var entrypoint string

		entrypoint = "."

		if 0 < len(args) {

			entrypoint = args[0]
		}

		/*
		fmt.Println("hof gen:",
			entrypoint,
		)
		*/

		err := pkg.Do(entrypoint)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	},
}

func init() {
	RootCmd.AddCommand(GenCmd)
}

func init() {
	// add sub-commands to this command when present

}
