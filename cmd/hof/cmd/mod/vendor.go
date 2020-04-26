package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
)

var vendorLong = `make a vendored copy of dependencies`

func VendorRun(args []string) (err error) {

	err = lib.ProcessLangs("vendor", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var VendorCmd = &cobra.Command{

	Use: "vendor [langs...]",

	Short: "make a vendored copy of dependencies",

	Long: vendorLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = VendorRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
