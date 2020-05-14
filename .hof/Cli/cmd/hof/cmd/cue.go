package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
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

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

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
