package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var bisectLong = `Use binary search to find the commit that introduced a bug`

func BisectRun(args []string) (err error) {

	return err
}

var BisectCmd = &cobra.Command{

	Use: "bisect",

	Short: "Use binary search to find the commit that introduced a bug",

	Long: bisectLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = BisectRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
