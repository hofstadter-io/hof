package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
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
