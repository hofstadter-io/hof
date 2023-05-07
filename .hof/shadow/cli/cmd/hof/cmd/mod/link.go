package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var linkLong = `symlink dependencies to cue.mod/pkg`

func LinkRun(args []string) (err error) {

	err = mod.Link(flags.RootPflags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var LinkCmd = &cobra.Command{

	Use: "link",

	Short: "symlink dependencies to cue.mod/pkg",

	Long: linkLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = LinkRun(args)
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

	ohelp := LinkCmd.HelpFunc()
	ousage := LinkCmd.UsageFunc()

	help := func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath() + " help")

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
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		return usage(cmd)
	}
	LinkCmd.SetHelpFunc(thelp)
	LinkCmd.SetUsageFunc(tusage)

}
