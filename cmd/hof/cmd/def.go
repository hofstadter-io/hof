package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/cuecmd"
)

var defLong = `print consolidated CUE definitions`

func init() {

	flags.SetupDefFlags(DefCmd.Flags(), &(flags.DefFlags))

}

func DefRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = cuecmd.Def(args, flags.RootPflags, flags.DefFlags)

	return err
}

var DefCmd = &cobra.Command{

	Use: "def",

	Short: "print consolidated CUE definitions",

	Long: defLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = DefRun(args)
		if err != nil {
			// fmt.Println(err)
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := DefCmd.HelpFunc()
	ousage := DefCmd.UsageFunc()

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
	DefCmd.SetHelpFunc(thelp)
	DefCmd.SetUsageFunc(tusage)

}
