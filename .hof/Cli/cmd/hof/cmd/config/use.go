package cmdconfig

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var useLong = `set the default configuration`

func UseRun(args []string) (err error) {

	return err
}

var UseCmd = &cobra.Command{

	Use: "use",

	Short: "set the default configuration",

	Long: useLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = UseRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
