package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/cuecmd"
)

var exportLong = `output data in a standard format`

func init() {
	ExportCmd.Flags().StringArrayVarP(&(flags.ExportFlags.Expression), "expression", "e", nil, "evaluate these expressions only")
	ExportCmd.Flags().BoolVarP(&(flags.ExportFlags.List), "list", "", false, "concatenate multiple objects into a list")
	ExportCmd.Flags().StringVarP(&(flags.ExportFlags.Out), "out", "", "", "output data format, when detection does not work")
	ExportCmd.Flags().StringVarP(&(flags.ExportFlags.Outfile), "outfile", "o", "", "filename or - for stdout with optional file prefix")
	ExportCmd.Flags().BoolVarP(&(flags.ExportFlags.Escape), "escape", "", false, "use HTLM escaping")
	ExportCmd.Flags().BoolVarP(&(flags.ExportFlags.Comments), "comments", "C", false, "include comments in output")
}

func ExportRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = cuecmd.Export(args, flags.RootPflags, flags.ExportFlags)

	return err
}

var ExportCmd = &cobra.Command{

	Use: "export",

	Short: "output data in a standard format",

	Long: exportLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = ExportRun(args)
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

	ohelp := ExportCmd.HelpFunc()
	ousage := ExportCmd.UsageFunc()

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
	ExportCmd.SetHelpFunc(thelp)
	ExportCmd.SetUsageFunc(tusage)

}
