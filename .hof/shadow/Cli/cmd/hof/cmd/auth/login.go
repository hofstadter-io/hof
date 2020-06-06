package cmdauth

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var loginLong = `login to an account, provider, system, or url`

func LoginRun(where string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var LoginCmd = &cobra.Command{

	Use: "login <where>",

	Short: "login to an account, provider, system, or url",

	Long: loginLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		var where string

		if 0 < len(args) {

			where = args[0]

		}

		err = LoginRun(where)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := LoginCmd.HelpFunc()
	usage := LoginCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "", 0)
		return usage(cmd)
	}
	LoginCmd.SetHelpFunc(thelp)
	LoginCmd.SetUsageFunc(tusage)

}
