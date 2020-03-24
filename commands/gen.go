package commands

import (
	"fmt"

	"github.com/spf13/cobra"
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

		// Default body

		fmt.Println("hof gen")

	},
}
