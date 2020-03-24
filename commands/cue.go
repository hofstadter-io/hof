package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cueLong = `Hof has cuelang embedded, so you can use hof anywhere you use cue`

var CueCmd = &cobra.Command{

	Use: "cue",

	Aliases: []string{
		"c",
	},

	Short: "Call a cue command",

	Long: cueLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		fmt.Println("run: cue", args)

	},
}
