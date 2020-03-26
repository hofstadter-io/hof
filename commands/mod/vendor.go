package mod

import (

	// hello... something might need to go here

	"github.com/spf13/cobra"

	"fmt"

	"os"

	"github.com/hofstadter-io/mvs/lib"
)

var vendorLong = `make a vendored copy of dependencies`

var VendorCmd = &cobra.Command{

	Use: "vendor [langs...]",

	Short: "make a vendored copy of dependencies",

	Long: vendorLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		err := lib.ProcessLangs("vendor", args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
