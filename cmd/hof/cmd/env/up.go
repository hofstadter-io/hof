package cmdenv

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	libenvcmd "github.com/hofstadter-io/hof/lib/env/cmd"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var upLong = `starts target points in an environment, this is very similar to docker-compose or helm locally`

func UpRun(args []string) (err error) {

	err = libenvcmd.Up(args, flags.RootPflags, flags.EnvPflags)

	return err
}

var UpCmd = &cobra.Command{

	Use: "up [...target] [% ...cue]",

	Short: "starts target points in an environment",

	Long: upLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = UpRun(args)
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

	ohelp := UpCmd.HelpFunc()
	ousage := UpCmd.UsageFunc()

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
	UpCmd.SetHelpFunc(thelp)
	UpCmd.SetUsageFunc(tusage)

}
