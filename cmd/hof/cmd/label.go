package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var labelLong = `manage resource labels`

func LabelRun(args []string) (err error) {

	return err
}

var LabelCmd = &cobra.Command{

	Use: "label",

	Short: "manage resource labels",

	Long: labelLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = LabelRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
