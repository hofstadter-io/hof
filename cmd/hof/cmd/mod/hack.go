package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
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
