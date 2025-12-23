package cmd

import (
	"fmt"
	"os"

	"github.com/dagger/dagger/analytics"
	"github.com/dagger/dagger/engine"
	enginetel "github.com/dagger/dagger/engine/telemetry"
	"github.com/spf13/cobra"

	cmdenv "github.com/hofstadter-io/hof/cmd/hof/cmd/env"
	"github.com/hofstadter-io/hof/lib/env/incept"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var envLong = `build, run, ship, and deploy environments (image, service, stack)`

func EnvPersistentPreRun(args []string) (err error) {
	// maybe the goes on the persistent preRunE
	workdir := "."
	workdir, err = incept.NormalizeWorkdir(workdir)
	if err != nil {
		return err
	}
	if err := os.Chdir(workdir); err != nil {
		return err
	}
	labels := enginetel.LoadDefaultLabels(workdir, engine.Version)
	t := analytics.New(analytics.DefaultConfig(labels))
	// cmd.SetContext(analytics.WithContext(cmd.Context(), t))
	cobra.OnFinalize(func() {
		t.Close()
	})

	// t.Capture(cmd.Context(), "cli_command", map[string]string{
	// 	"name": commandName(cmd),
	// })

	return err
}

var EnvCmd = &cobra.Command{

	Use: "env [args]",

	Short: "build, run, ship, and deploy environments (image, service, stack)",

	Long: envLong,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = EnvPersistentPreRun(args)
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

	ohelp := EnvCmd.HelpFunc()
	ousage := EnvCmd.UsageFunc()

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
	EnvCmd.SetHelpFunc(thelp)
	EnvCmd.SetUsageFunc(tusage)

	EnvCmd.AddCommand(cmdenv.BuildCmd)
	EnvCmd.AddCommand(cmdenv.InfoCmd)
	EnvCmd.AddCommand(cmdenv.ListCmd)
	EnvCmd.AddCommand(cmdenv.ImagesCmd)
	EnvCmd.AddCommand(cmdenv.PsCmd)
	EnvCmd.AddCommand(cmdenv.RunCmd)
	EnvCmd.AddCommand(cmdenv.UpCmd)
	EnvCmd.AddCommand(cmdenv.DownCmd)
	EnvCmd.AddCommand(cmdenv.TagCmd)
	EnvCmd.AddCommand(cmdenv.PushCmd)
	EnvCmd.AddCommand(cmdenv.PullCmd)
	EnvCmd.AddCommand(cmdenv.DeployCmd)

}
