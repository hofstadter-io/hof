package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var initLong = `create an empty workspace or initialize an existing directory to one

  module name or path should look like github.com/hofstadter-io/hof`

func InitRun(module string, name string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var InitCmd = &cobra.Command{

	Use: "init <module>",

	Short: "create an empty workspace or initialize an existing directory to one",

	Long: initLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'module'")
			cmd.Usage()
			os.Exit(1)
		}

		var module string

		if 0 < len(args) {

			module = args[0]

		}

		var name string

		if 1 < len(args) {

			name = args[1]

		}

		err = InitRun(module, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := InitCmd.HelpFunc()
	usage := InitCmd.UsageFunc()

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
	InitCmd.SetHelpFunc(thelp)
	InitCmd.SetUsageFunc(tusage)

}
