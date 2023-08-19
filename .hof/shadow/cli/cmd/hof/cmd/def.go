package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var defLong = `print consolidated CUE definitions`

func init() {

	DefCmd.Flags().StringSliceVarP(&(flags.DefFlags.Expression), "expression", "e", nil, "evaluate these expressions only")
	DefCmd.Flags().BoolVarP(&(flags.DefFlags.List), "list", "", false, "concatenate multiple objects into a list")
	DefCmd.Flags().StringVarP(&(flags.DefFlags.Out), "out", "", "", "output data format, when detection does not work")
	DefCmd.Flags().StringVarP(&(flags.DefFlags.Outfile), "outfile", "o", "", "filename or - for stdout with optional file prefix")
	DefCmd.Flags().BoolVarP(&(flags.DefFlags.InlineImports), "inline-imports", "", false, "expand references to non-core imports")
	DefCmd.Flags().BoolVarP(&(flags.DefFlags.Comments), "comments", "C", false, "include comments in output")
	DefCmd.Flags().BoolVarP(&(flags.DefFlags.Attributes), "attributes", "a", false, "diplay field attributes")
}

func DefRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

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
