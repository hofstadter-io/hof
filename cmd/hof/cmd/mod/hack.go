package cmdmod

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var hackLong = `dev command`

func HackRun(args []string) (err error) {

	err = lib.Hack("", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var HackCmd = &cobra.Command{

	Use: "hack",

	Hidden: true,

	Short: "dev command",

	Long: hackLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = HackRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
