package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testLong = `test your auth configuration, defaults to current`

var TestCmd = &cobra.Command{

	Use: "test [name]",

	Short: "test your auth configuration, defaults to current",

	Long: testLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		var name string

		if 0 < len(args) {

			name = args[0]

		}

		// Default body

		fmt.Println("hof config test", name)

	},
}
