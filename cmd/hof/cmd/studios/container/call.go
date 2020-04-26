package cmdcontainer

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var callLong = `Call your Studios container`

func CallRun(args []string) (err error) {

	return err
}

var CallCmd = &cobra.Command{

	Use: "call",

	Short: "Call a container",

	Long: callLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = CallRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
