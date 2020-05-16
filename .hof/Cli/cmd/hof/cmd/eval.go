package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var evalLong = `print consolidated definitions`

func EvalRun(args []string) (err error) {

	return err
}

var EvalCmd = &cobra.Command{

	Use: "eval",

	Short: "print consolidated definitions",

	Long: evalLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = EvalRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
