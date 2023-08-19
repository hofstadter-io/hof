package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/cuecmd"
)

var vetLong = `validate data with CUE`

func init() {

	VetCmd.Flags().StringSliceVarP(&(flags.VetFlags.Expression), "expression", "e", nil, "evaluate these expressions only")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.Extensions), "extensions", "x", false, "include hof extensions when evaluating CUE code")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.List), "list", "", false, "concatenate multiple objects into a list")
	VetCmd.Flags().StringVarP(&(flags.VetFlags.Out), "out", "", "", "output data format, when detection does not work")
	VetCmd.Flags().StringVarP(&(flags.VetFlags.Outfile), "outfile", "o", "", "filename or - for stdout with optional file prefix")
	VetCmd.Flags().StringVarP(&(flags.VetFlags.Schema), "schema", "d", "", "expression to select schema for evaluating values in non-CUE files")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.Concrete), "concrete", "c", false, "require the evaluation to be concrete")
}

func VetRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = cuecmd.Vet(args, flags.RootPflags, flags.VetFlags)

	return err
}

var VetCmd = &cobra.Command{

	Use: "vet",

	Short: "validate data with CUE",

	Long: vetLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = VetRun(args)
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

	ohelp := VetCmd.HelpFunc()
	ousage := VetCmd.UsageFunc()

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
	VetCmd.SetHelpFunc(thelp)
	VetCmd.SetUsageFunc(tusage)

}
