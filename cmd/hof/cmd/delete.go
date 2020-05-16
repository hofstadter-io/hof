package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var deleteLong = `delete resources`

func DeleteRun(args []string) (err error) {

	return err
}

var DeleteCmd = &cobra.Command{

	Use: "delete",

	Short: "delete resources",

	Long: deleteLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = DeleteRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
