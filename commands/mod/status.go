package mod

import (

	// hello... something might need to go here

	"github.com/spf13/cobra"

	"fmt"

	"os"

	"github.com/hofstadter-io/mvs/lib"
)

var statusLong = `print module dependencies status`

var StatusCmd = &cobra.Command{

	Use: "status",

	Short: "print module dependencies status",

	Long: statusLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		err := lib.ProcessLangs("status", args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
