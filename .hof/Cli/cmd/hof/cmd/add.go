package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var addLong = `add dependencies and new components to the current module or workspace`

func AddRun(args []string) (err error) {

	return err
}

var AddCmd = &cobra.Command{

	Use: "add",

	Short: "add dependencies and new components to the current module or workspace",

	Long: addLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = AddRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
