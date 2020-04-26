package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cueLong = `Hof has cuelang embedded, so you can use hof anywhere you use cue`

func CueRun(args []string) (err error) {

	return err
}

var CueCmd = &cobra.Command{

	Use: "cue",

	Aliases: []string{
		"c",
	},

	Short: "Call a cue command",

	Long: cueLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = CueRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
