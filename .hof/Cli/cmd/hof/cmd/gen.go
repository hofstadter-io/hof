package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var genLong = `  generate all the things, from code to data to config...`

var (
	GenGeneratorFlag []string
)

func init() {

	GenCmd.Flags().StringSliceVarP(&GenGeneratorFlag, "generator", "g", nil, "Generators to run")
}

func init() {

}

func GenRun(args []string) (err error) {

	return err
}

var GenCmd = &cobra.Command{

	Use: "gen [files...]",

	Aliases: []string{
		"g",
	},

	Short: "generate code, data, and config",

	Long: genLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

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
