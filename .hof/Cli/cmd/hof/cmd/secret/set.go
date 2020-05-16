package cmdsecret

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var setLong = `set secret value(s)`

func SetRun(args []string) (err error) {

	return err
}

var SetCmd = &cobra.Command{

	Use: "set",

	Short: "set secret value(s)",

	Long: setLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = SetRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
