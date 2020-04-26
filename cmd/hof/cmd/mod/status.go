package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
)

var statusLong = `print module dependencies status`

func StatusRun(args []string) (err error) {

	err = lib.ProcessLangs("status", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var StatusCmd = &cobra.Command{

	Use: "status",

	Short: "print module dependencies status",

	Long: statusLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = StatusRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
