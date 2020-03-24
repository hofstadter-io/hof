package mod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
)

var graphLong = `print module requirement graph`

var GraphCmd = &cobra.Command{

	Use: "graph",

	Short: "print module requirement graph",

	Long: graphLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		err := lib.ProcessLangs("graph", args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
