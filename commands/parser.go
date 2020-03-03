package commands

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/lang"
)

// Tool:   hof
// Name:   gen
// Usage:  gen <entrypoint>
// Parent: hof

var ParserLong = `Parsererate a project starting from the file.`

var ParserCmd = &cobra.Command{

	Use: "parser <entrypoint>",

	Short: "Parsererate a project.",

	Long: ParserLong,

	Run: func(cmd *cobra.Command, args []string) {

		err := lang.Parser(args)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	},
}

func init() {
	RootCmd.AddCommand(ParserCmd)
}

func init() {
	// add sub-commands to this command when present

}
