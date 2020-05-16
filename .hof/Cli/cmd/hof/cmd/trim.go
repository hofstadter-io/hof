package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var trimLong = `cleanup code, configuration, and more`

func TrimRun(args []string) (err error) {

	return err
}

var TrimCmd = &cobra.Command{

	Use: "trim",

	Short: "cleanup code, configuration, and more",

	Long: trimLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = TrimRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
