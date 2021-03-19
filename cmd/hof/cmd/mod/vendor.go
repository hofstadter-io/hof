package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"
)

var vendorLong = `make a vendored copy of dependencies`

func VendorRun(args []string) (err error) {

	err = mod.ProcessLangs("vendor", args)
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

	PreRun: func(cmd *cobra.Command, args []string) {

	},

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

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := VendorCmd.HelpFunc()
	ousage := VendorCmd.UsageFunc()
	help := func(cmd *cobra.Command, args []string) {
		if extra(cmd) {
			return
		}
		ohelp(cmd, args)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		return ousage(cmd)
	}

	VendorCmd.SetHelpFunc(help)
	VendorCmd.SetUsageFunc(usage)

}
