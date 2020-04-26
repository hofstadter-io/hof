package cmdconfig

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var testLong = `test your auth configuration, defaults to current`

func TestRun(name string) (err error) {

	return err
}

var TestCmd = &cobra.Command{

	Use: "test [name]",

	Short: "test your auth configuration, defaults to current",

	Long: testLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		var name string

		if 0 < len(args) {

			name = args[0]

		}

		err = TestRun(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
