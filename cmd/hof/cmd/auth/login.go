package cmdauth

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

		err = LoginRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
