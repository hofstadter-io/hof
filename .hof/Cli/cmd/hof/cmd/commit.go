package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var commitLong = `Record changes to the repository`

func CommitRun(args []string) (err error) {

	return err
}

var CommitCmd = &cobra.Command{

	Use: "commit",

	Short: "Record changes to the repository",

	Long: commitLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = CommitRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
