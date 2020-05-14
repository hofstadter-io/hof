package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var newLong = `create a new project or subcomponent or files, depending on the context`

func NewRun(args []string) (err error) {

	return err
}

var NewCmd = &cobra.Command{

	Use: "new",

	Short: "create a new project or subcomponent or files",

	Long: newLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = NewRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
