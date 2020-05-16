package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var applyLong = `apply resource configuration`

func ApplyRun(args []string) (err error) {

	return err
}

var ApplyCmd = &cobra.Command{

	Use: "apply",

	Short: "apply resource configuration",

	Long: applyLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ApplyRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
