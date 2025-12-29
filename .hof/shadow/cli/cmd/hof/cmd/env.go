package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/env"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var envLong = `build, run, ship, and deploy environments (image, service, stack)`

func init() {

	flags.SetupEnvPflags(EnvCmd.PersistentFlags(), &(flags.EnvPflags))

}

func EnvRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var EnvCmd = &cobra.Command{

	Use: "env [args]",

	Short: "build, run, ship, and deploy environments (image, service, stack)",

	Long: envLong,

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		glob := toComplete + "*"
		matches, _ := filepath.Glob(glob)
		return matches, cobra.ShellCompDirectiveDefault
	},

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = EnvRun(args)
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
	EnvCmd.AddCommand(cmdenv.ExportCmd)
	EnvCmd.AddCommand(cmdenv.GetCmd)
	EnvCmd.AddCommand(cmdenv.ListCmd)
	EnvCmd.AddCommand(cmdenv.RunCmd)
	EnvCmd.AddCommand(cmdenv.UpCmd)
	EnvCmd.AddCommand(cmdenv.PublishCmd)

}
