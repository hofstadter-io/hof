package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var statusLong = `Show workspace information and status`

func StatusRun(args []string) (err error) {

	return err
}

var StatusCmd = &cobra.Command{

	Use: "status",

	Short: "Show workspace information and status",

	Long: statusLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = StatusRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
