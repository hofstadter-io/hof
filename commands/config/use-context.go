package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/config"
)

var UseContextLong = `Sets the context as the default`

var UseContextCmd = &cobra.Command{

	Use: "use-context <context-name>",

	Aliases: []string{
		"use",
	},

	Short: "Use a context configuration",

	Long: UseContextLong,

	Run: func(cmd *cobra.Command, args []string) {

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'context'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var context string
		if 0 < len(args) {
			context = args[0]
		}

		/*
			fmt.Println("hof config use-context:",
				context,
			)
		*/

		config.UseContext(context)

	},
}
