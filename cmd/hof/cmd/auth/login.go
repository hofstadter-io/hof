package cmdauth

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var loginLong = `login to an account`

func LoginRun(args []string) (err error) {

	return err
}

var LoginCmd = &cobra.Command{

	Use: "login",

	Short: "login to an account",

	Long: loginLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = LoginRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
