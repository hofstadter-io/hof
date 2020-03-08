package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/old-lang"
)

var GenLong = `Generate a project starting from the file.`

var GenCmd = &cobra.Command{

	Use: "gen <entrypoint>",

	Aliases: []string{
		"g",
	},

	Short: "Generate a project.",

	Long: GenLong,

	Run: func(cmd *cobra.Command, args []string) {

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

		// err := lang.Gen(entrypoint)
		err := lang.Eval(entrypoint)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	},
}

func init() {
	RootCmd.AddCommand(GenCmd)
}
