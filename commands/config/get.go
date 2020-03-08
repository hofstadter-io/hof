package config

import (
	// "fmt"

	// custom imports

	// infered imports

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/config"
)

// Tool:   hof
// Name:   get
// Usage:  get [all|<context-name>]
// Parent: config

var GetLong = `Get your configuration`

var GetCmd = &cobra.Command{

	Use: "get [all|<context-name>]",

	Aliases: []string{
		"view",
	},

	Short: "Get Hofsadter configuration(s)",

	Long: GetLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In getCmd", "args", args)
		// Argument Parsing
		// [0]name:   context
		//     help:
		//     req'd:

		var context string

		if 0 < len(args) {

			context = args[0]
		}

		/*
		fmt.Println("hof config get:",
			context,
		)
		*/

		config.GetContext(context)
	},
}

func init() {
	// add sub-commands to this command when present

}
