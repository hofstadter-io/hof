package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib"
)

var genLong = `  generate all the things, from code to data to config...`

func GenRun(args []string) (err error) {

	return lib.Gen([]string{}, []string{}, "")

	return err
}

var GenCmd = &cobra.Command{

	Use: "gen [files...]",

	Aliases: []string{
		"g",
	},

	Short: "generate code, data, and config",

	Long: genLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = GenRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
