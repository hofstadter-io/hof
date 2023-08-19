package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var evalLong = `evaluate and print CUE configuration`

func init() {

	EvalCmd.Flags().StringSliceVarP(&(flags.EvalFlags.Expression), "expression", "e", nil, "evaluate these expressions only")
	EvalCmd.Flags().BoolVarP(&(flags.EvalFlags.Extensions), "extensions", "x", false, "include hof extensions when evaluating CUE code")
	EvalCmd.Flags().BoolVarP(&(flags.EvalFlags.List), "list", "", false, "concatenate multiple objects into a list")
	EvalCmd.Flags().StringVarP(&(flags.EvalFlags.Out), "out", "", "", "output data format, when detection does not work")
	EvalCmd.Flags().StringVarP(&(flags.EvalFlags.Outfile), "outfile", "o", "", "filename or - for stdout with optional file prefix")
	EvalCmd.Flags().StringVarP(&(flags.EvalFlags.Schema), "schema", "d", "", "expression to select schema for evaluating values in non-CUE files")
	EvalCmd.Flags().BoolVarP(&(flags.EvalFlags.All), "all", "a", false, "show optional and hidden fields")
	EvalCmd.Flags().BoolVarP(&(flags.EvalFlags.Concrete), "concrete", "c", false, "require the evaluation to be concrete")
	EvalCmd.Flags().BoolVarP(&(flags.EvalFlags.Attributes), "attributes", "A", false, "diplay field attributes")
	EvalCmd.Flags().BoolVarP(&(flags.EvalFlags.Hidden), "hidden", "H", false, "diplay hidden fields")
	EvalCmd.Flags().BoolVarP(&(flags.EvalFlags.Optional), "optional", "O", false, "diplay optional fields")
}

func EvalRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var EvalCmd = &cobra.Command{

	Use: "eval",

	Short: "evaluate and print CUE configuration",

	Long: evalLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = EvalRun(args)
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

	ohelp := EvalCmd.HelpFunc()
	ousage := EvalCmd.UsageFunc()

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
	EvalCmd.SetHelpFunc(thelp)
	EvalCmd.SetUsageFunc(tusage)

}
