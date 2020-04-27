package cmdauth

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/pflags"
)

var loginLong = `login to an account`

func LoginRun(args []string) (err error) {

	return err
}

var LoginCmd = &cobra.Command{

	Use: "login",

	Short: "login to an account",

	Long: loginLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		fmt.Printf("login: %q\n", pflags.RootConfigPflag)

		err = LoginRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
