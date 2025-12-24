package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/env"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var envLong = `build, run, ship, and deploy environments (image, service, stack)`

func init() {

	flags.SetupEnvPflags(EnvCmd.PersistentFlags(), &(flags.EnvPflags))

}

var EnvCmd = &cobra.Command{

	Use: "env [args]",

	Short: "build, run, ship, and deploy environments (image, service, stack)",

	Long: envLong,
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
	EnvCmd.AddCommand(cmdenv.ExportCmd)
	EnvCmd.AddCommand(cmdenv.InfoCmd)
	EnvCmd.AddCommand(cmdenv.ListCmd)
	EnvCmd.AddCommand(cmdenv.ImagesCmd)
	EnvCmd.AddCommand(cmdenv.PsCmd)
	EnvCmd.AddCommand(cmdenv.RunCmd)
	EnvCmd.AddCommand(cmdenv.UpCmd)
	EnvCmd.AddCommand(cmdenv.DownCmd)
	EnvCmd.AddCommand(cmdenv.TagCmd)
	EnvCmd.AddCommand(cmdenv.PublishCmd)
	EnvCmd.AddCommand(cmdenv.PushCmd)
	EnvCmd.AddCommand(cmdenv.PullCmd)
	EnvCmd.AddCommand(cmdenv.DeployCmd)

}
