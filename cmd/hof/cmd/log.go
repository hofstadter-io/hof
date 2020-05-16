package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var logLong = `Show workspace logs and history`

func LogRun(args []string) (err error) {

	return err
}

var LogCmd = &cobra.Command{

	Use: "log",

	Short: "Show workspace logs and history",

	Long: logLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = LogRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
