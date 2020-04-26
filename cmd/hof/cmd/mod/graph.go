package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
)

var graphLong = `print module requirement graph`

func GraphRun(args []string) (err error) {

	err = lib.ProcessLangs("graph", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var GraphCmd = &cobra.Command{

	Use: "graph",

	Short: "print module requirement graph",

	Long: graphLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = GraphRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
