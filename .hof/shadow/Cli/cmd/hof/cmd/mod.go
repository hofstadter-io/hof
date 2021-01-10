package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/mod"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var modLong = `The mod subcmd is a polyglot dependency management tool based on go mods.

mod file format:

  module = "<module path>"

  <name> = "version"

  require (
    ...
  )

  replace <module path> => <local path>
  ...`

func ModPersistentPreRun(args []string) (err error) {

	mod.InitLangs()

	return err
}

var ModCmd = &cobra.Command{

	Use: "mod",

	Aliases: []string{
		"m",
	},

	Short: "mod subcmd is a polyglot dependency management tool based on go mods",

	Long: modLong,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ModPersistentPreRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := ModCmd.HelpFunc()
	ousage := ModCmd.UsageFunc()
	help := func(cmd *cobra.Command, args []string) {
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
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	ModCmd.SetHelpFunc(thelp)
	ModCmd.SetUsageFunc(tusage)

	ModCmd.AddCommand(cmdmod.InfoCmd)
	ModCmd.AddCommand(cmdmod.InitCmd)
	ModCmd.AddCommand(cmdmod.VendorCmd)

}
