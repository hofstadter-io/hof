package cmdconfig

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var getLong = `print a configuration`

func GetRun(name string) (err error) {

	return err
}

var GetCmd = &cobra.Command{

	Use: "get",

	Short: "print a configuration",

	Long: getLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Name'")
			cmd.Usage()
			os.Exit(1)
		}

		var name string

		if 0 < len(args) {

			name = args[0]

		}

		err = GetRun(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
