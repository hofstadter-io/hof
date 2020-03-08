package config

import (
	"fmt"

	// custom imports

	// infered imports
	"os"

	"github.com/spf13/cobra"
	"github.com/hofstadter-io/hof/lib/config"
)

// Tool:   hof
// Name:   use-context
// Usage:  use-context <context-name>
// Parent: config

var UseContextLong = `Sets the context as the default`

var UseContextCmd = &cobra.Command{

	Use: "use-context <context-name>",

	Aliases: []string{
		"use",
	},

	Short: "Use a context configuration",

	Long: UseContextLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In use-contextCmd", "args", args)
		// Argument Parsing
		// [0]name:   context
		//     help:
		//     req'd:  true
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

func init() {
	// add sub-commands to this command when present

}
