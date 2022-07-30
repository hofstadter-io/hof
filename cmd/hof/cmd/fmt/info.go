package cmdfmt

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	hfmt "github.com/hofstadter-io/hof/lib/fmt"
)

var infoLong = `get formatter info`

func InfoRun(formatter string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = hfmt.Info(formatter)

	return err
}

var InfoCmd = &cobra.Command{

	Use: "info",

	Short: "get formatter info",

	Long: infoLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		var formatter string

		if 0 < len(args) {

			formatter = args[0]

		}

		err = InfoRun(formatter)
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

	ohelp := InfoCmd.HelpFunc()
	ousage := InfoCmd.UsageFunc()
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

	InfoCmd.SetHelpFunc(help)
	InfoCmd.SetUsageFunc(usage)

}
