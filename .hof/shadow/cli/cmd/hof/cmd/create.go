package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var createLong = `hof create enables you to easily bootstrap
code for full projects, components, and more.

Examples can be found in the documentation:

  https://docs.hofstadter.io/hof-create/

By adding one config file and templates to your repo
your users can quickly bootstrap full applications,
tooling configuration, and other code using your project.
Share consistent scaffolding, configurable to users.

Any hof generator can also support the create command
and most choose to bootstrap a generator at minimum.
This means you get all the same benefits from
hof's code generation engine, turning your
bootstrapped code into a living template.

Run create from any git repo and any ref

  hof create github.com/username/repo@v1.2.3
  hof create github.com/username/repo@a1b2c3f
  hof create github.com/username/repo@latest

-I supplies inputs as key/value pairs or from a file
when no flag is supplied, an interactive prompt is used

  hof create github.com/username/repo@v1.2.3 \
    -I name=foo -I val=bar \
    -I @inputs.cue

You can also reference local generators by their cue inputs.
This local lookup is indicated by ./ or ../ starting a path.
Use this mode when developing and testing locally.

  hof create ../my-gen`

func init() {

	CreateCmd.Flags().StringSliceVarP(&(flags.CreateFlags.Input), "input", "I", nil, "inputs to the create module")
	CreateCmd.Flags().StringSliceVarP(&(flags.CreateFlags.Generator), "generator", "G", nil, "generator tags to run, default is all")
	CreateCmd.Flags().StringVarP(&(flags.CreateFlags.Outdir), "outdir", "O", "", "base directory to write all output to")
}

func CreateRun(module string, extra []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var CreateCmd = &cobra.Command{

	Use: "create <module path> [extra args]",

	Short: "starter kits or blueprints from any git repo",

	Long: createLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'module'")
			cmd.Usage()
			os.Exit(1)
		}

		var module string

		if 0 < len(args) {

			module = args[0]

		}

		var extra []string

		if 1 < len(args) {

			extra = args[1:]

		}

		err = CreateRun(module, extra)
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

	ohelp := CreateCmd.HelpFunc()
	ousage := CreateCmd.UsageFunc()

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
	CreateCmd.SetHelpFunc(thelp)
	CreateCmd.SetUsageFunc(tusage)

}
