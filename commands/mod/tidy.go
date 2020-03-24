package mod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
)

var tidyLong = `add missinad and remove unused modules`

var TidyCmd = &cobra.Command{

	Use: "tidy [langs...]",

	Short: "add missinad and remove unused modules",

	Long: tidyLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		err := lib.ProcessLangs("tidy", args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
