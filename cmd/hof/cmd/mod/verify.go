package cmdmod

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var verifyLong = `verify dependencies have expected content`

func VerifyRun(args []string) (err error) {

	err = lib.ProcessLangs("verify", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var VerifyCmd = &cobra.Command{

	Use: "verify [langs...]",

	Short: "verify dependencies have expected content",

	Long: verifyLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = VerifyRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
