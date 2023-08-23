package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var vetLong = `validate data with CUE`

func init() {

	VetCmd.Flags().StringSliceVarP(&(flags.VetFlags.Expression), "expression", "e", nil, "evaluate these expressions only")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.List), "list", "", false, "concatenate multiple objects into a list")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.Simplify), "simplify", "", false, "simplify CUE statements where possible")
	VetCmd.Flags().StringVarP(&(flags.VetFlags.Out), "out", "", "", "output data format, when detection does not work")
	VetCmd.Flags().StringVarP(&(flags.VetFlags.Outfile), "outfile", "o", "", "filename or - for stdout with optional file prefix")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.Concrete), "concrete", "c", false, "require the evaluation to be concrete")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.Comments), "comments", "C", false, "include comments in output")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.Attributes), "attributes", "A", false, "display field attributes")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.Definitions), "definitions", "S", true, "display defintions")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.Hidden), "hidden", "H", false, "display hidden fields")
	VetCmd.Flags().BoolVarP(&(flags.VetFlags.Optional), "optional", "O", false, "display optional fields")
}

func VetRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

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
