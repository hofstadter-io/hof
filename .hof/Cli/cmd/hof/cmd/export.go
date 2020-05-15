package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var exportLong = `export your data model to various formats`

func ExportRun(args []string) (err error) {

	return err
}

var ExportCmd = &cobra.Command{

	Use: "export",

	Short: "export your data model to various formats",

	Long: exportLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ExportRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
