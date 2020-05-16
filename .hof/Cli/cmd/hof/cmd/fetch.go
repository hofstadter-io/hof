package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var fetchLong = `Download objects and refs from another repository`

func FetchRun(args []string) (err error) {

	return err
}

var FetchCmd = &cobra.Command{

	Use: "fetch",

	Short: "Download objects and refs from another repository",

	Long: fetchLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = FetchRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
