package config

import (
	// "fmt"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/config"
)

var GetLong = `Get your configuration`

var GetCmd = &cobra.Command{

	Use: "get [all|<context-name>]",

	Aliases: []string{
		"view",
	},

	Short: "Get Hofsadter configuration(s)",

	Long: GetLong,

	Run: func(cmd *cobra.Command, args []string) {

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
