package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var includeLong = `Include changes into the changeset`

func IncludeRun(args []string) (err error) {

	return err
}

var IncludeCmd = &cobra.Command{

	Use: "include",

	Short: "Include changes into the changeset",

	Long: includeLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = IncludeRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
