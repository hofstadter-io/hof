package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
	dcmd "github.com/hofstadter-io/hof/lib/dagger/cmd"
)

var daggerooLong = `dagger run helper for the extension server`

func init() {

	flags.SetupDaggerooFlags(DaggerooCmd.Flags(), &(flags.DaggerooFlags))

}

func DaggerooRun(id string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	return dcmd.Run(id, flags.RootPflags, flags.DaggerooFlags)
}

var DaggerooCmd = &cobra.Command{

	Use: "daggeroo [args]",

	Short: "dagger run helper for the extension server",

	Long: daggerooLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'id'")
			cmd.Usage()
			os.Exit(1)
		}

		var id string

		if 0 < len(args) {

			id = args[0]

		}

		err = DaggerooRun(id)
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

	ohelp := DaggerooCmd.HelpFunc()
	ousage := DaggerooCmd.UsageFunc()

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
	DaggerooCmd.SetHelpFunc(thelp)
	DaggerooCmd.SetUsageFunc(tusage)

}
