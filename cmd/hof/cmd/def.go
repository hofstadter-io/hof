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

	DefCmd.Flags().StringArrayVarP(&(flags.DefFlags.Expression), "expression", "e", nil, "evaluate these expressions only")
	DefCmd.Flags().BoolVarP(&(flags.DefFlags.Extensions), "extensions", "x", false, "include hof extensions when evaluating CUE code")
	DefCmd.Flags().BoolVarP(&(flags.DefFlags.List), "list", "", false, "concatenate multiple objects into a list")
	DefCmd.Flags().StringVarP(&(flags.DefFlags.Out), "out", "", "", "output data format, when detection does not work")
	DefCmd.Flags().StringVarP(&(flags.DefFlags.Outfile), "outfile", "o", "", "filename or - for stdout with optional file prefix")
	DefCmd.Flags().StringVarP(&(flags.DefFlags.Schema), "schema", "d", "", "expression to select schema for evaluating values in non-CUE files")
	DefCmd.Flags().BoolVarP(&(flags.DefFlags.InlineImports), "inline-imports", "", false, "expand references to non-core imports")
	DefCmd.Flags().BoolVarP(&(flags.DefFlags.Attributes), "attributes", "a", false, "diplay field attributes")
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
			fmt.Println(err)
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
