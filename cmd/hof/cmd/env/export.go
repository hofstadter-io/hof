package cmdenv

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/cmd/hof/ga"
	"github.com/hofstadter-io/hof/lib/env/cmd"
)

var exportLong = `export an environment into local container runtime`

func init() {

	flags.SetupEnv__ExportFlags(ExportCmd.Flags(), &(flags.Env__ExportFlags))

}

func ExportRun(args []string) (err error) {

	return cmd.Export(args, flags.RootPflags, flags.EnvPflags, flags.Env__ExportFlags)

	// // you can safely comment this print out
	// fmt.Println("not implemented")

	// return err
}

var ExportCmd = &cobra.Command{

	Use: "export [target...]",

	Short: "export an environment into local container runtime",

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
