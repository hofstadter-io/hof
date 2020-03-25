package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib"
)

var genLong = `  generate all the things, from code to data to config...`

var GenCmd = &cobra.Command{

	Use: "gen [files...]",

	Aliases: []string{
		"g",
	},

	Short: "generate code, data, and config",

	Long: genLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		msg, err := lib.Gen(args, []string{}, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(msg)

	},
}
