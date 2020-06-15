package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var verifyLong = `verify dependencies have expected content`

func VerifyRun(args []string) (err error) {

	err = mod.ProcessLangs("verify", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var VerifyCmd = &cobra.Command{

	Use: "verify [langs...]",

	Short: "verify dependencies have expected content",

	Long: verifyLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = VerifyRun(args)
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

	ohelp := VerifyCmd.HelpFunc()
	ousage := VerifyCmd.UsageFunc()
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

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	VerifyCmd.SetHelpFunc(thelp)
	VerifyCmd.SetUsageFunc(tusage)

}
