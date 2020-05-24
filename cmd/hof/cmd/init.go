package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/workspace"
)

var initLong = `Create an empty workspace or initialize an existing directory to one`

func InitRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = workspace.InitRunFromArgs(args)

	return err
}

var InitCmd = &cobra.Command{

	Use: "init",

	Short: "Create an empty workspace or initialize an existing directory to one",

	Long: initLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = InitRun(args)
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
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "<omit>", 0)
		return usage(cmd)
	}
	InitCmd.SetHelpFunc(thelp)
	InitCmd.SetUsageFunc(tusage)

}
