package cmdenv

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/cmd/hof/ga"
	"github.com/hofstadter-io/hof/lib/env/cmd"
)

var syncLong = `sync target points in an environment, making sure they are ready to go, no matter the type`

func SyncRun(args []string) (err error) {

	return cmd.Sync(args, flags.RootPflags, flags.EnvPflags)

	// you can safely comment this print out
	// fmt.Println("not implemented")

	// return err
}

var SyncCmd = &cobra.Command{

	Use: "sync [...target] [% ...cue]",

	Short: "sync target points in an environment",

	Long: syncLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = SyncRun(args)
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

	ohelp := SyncCmd.HelpFunc()
	ousage := SyncCmd.UsageFunc()

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
	SyncCmd.SetHelpFunc(thelp)
	SyncCmd.SetUsageFunc(tusage)

}
