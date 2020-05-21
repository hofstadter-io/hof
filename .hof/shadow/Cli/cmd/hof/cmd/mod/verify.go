package cmdmod

import (
	"fmt"
	"os"

	"strings"

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

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

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

	help := VerifyCmd.HelpFunc()
	usage := VerifyCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "<omit>", 0)
		return usage(cmd)
	}
	VerifyCmd.SetHelpFunc(thelp)
	VerifyCmd.SetUsageFunc(tusage)

}
