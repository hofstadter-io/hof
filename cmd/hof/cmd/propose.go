package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var proposeLong = `Propose to include your changeset in a remote repository`

func ProposeRun(args []string) (err error) {

	return err
}

var ProposeCmd = &cobra.Command{

	Use: "propose",

	Short: "Propose to include your changeset in a remote repository",

	Long: proposeLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ProposeRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
