package config

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

var getLong = `print a configuration`

var GetCmd = &cobra.Command{

	Use: "get",

	Short: "print a configuration",

	Long: getLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Name'")
			cmd.Usage()
			os.Exit(1)
		}

		var name string

		if 0 < len(args) {

			name = args[0]

		}

		// Default body

		fmt.Println("hof config get", name)

	},
}
