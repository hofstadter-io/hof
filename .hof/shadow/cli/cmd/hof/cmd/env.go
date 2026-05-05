package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/runtime"

	libenvcmd "github.com/hofstadter-io/hof/lib/env/cmd"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/env"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var envLong = `build, run, ship, and deploy environments (image, service, stack)

'veg env' looks for custom commands and treats them equally to builtin commands.
All commands have well known behaviors, depending on the $kind of a target.
'veg env list' and 'hof env sync' work with all '$kind's.
Most commands only work with a subset that makes sense or is explicit.
See their help text to learn more.

Note, hof -> veg in terms of name, I have started with this command, and also the agentic stuff
The big move is coming soon...

## Examples

# [...targets], or "points", are selected via args and flags
#   list allows you to explore that space without syncing or triggering evaluation
veg env list|info ['^name$'] [-K '^kind$'] [-P '^path$'] [-S name|kind|path]

# sync, evaluates the DAGs, but doesn't export or server anything, steps are run
#   it can act as "real dry run" compared to the --dry-run flag for other commands
#   without any args or flags, sync acts as a big test as well
veg env sync [...targets] [...flags]

# run, creates an interactive session and binds and dependent services
#   this is closest to docker run or kubectl exec
veg env run [...target] [...flags]

# run, launches an services or stacks, similar to compose and helm
veg env up [...target] [...flags]

# export artifacts, local or remote, object storage and registries
veg env export -P release -T v0.4.3 -t dest=./release

# make your own commands and flags, designed for your workflows
#   this is closest to Makefiles or package.json scripts
#   define similar commands with the power of CUE and Dagger
veg env [init, test, lint, ci, publish, deploy, ...]
veg env ... -t env=stg -t stack=app -t branch=main

## Important References

./schemas/env    # the CUE schemas for what you can do in veg/env
./catalogs/env   # reusable CUE for all sorts of things from small to big
./examples/env   # simple and complex examples for you to play and fork
`

func init() {

	flags.SetupEnvPflags(EnvCmd.PersistentFlags(), &(flags.EnvPflags))

}

func EnvPersistentPreRun(args []string) (err error) {

	err = runtime.EnsureInfra()

	return err
}

func EnvRun(args []string) (err error) {

	err = libenvcmd.Env(args, flags.RootPflags, flags.EnvPflags)

	return err
}

var EnvCmd = &cobra.Command{

	Use: "env [...target] [% ...cue]",

	Short: "build, run, ship, and deploy environments (image, service, stack)",

	Long: envLong,

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		glob := toComplete + "*"
		matches, _ := filepath.Glob(glob)
		return matches, cobra.ShellCompDirectiveDefault
	},

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = EnvPersistentPreRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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

	EnvCmd.AddCommand(cmdenv.SyncCmd)
	EnvCmd.AddCommand(cmdenv.ExportCmd)
	EnvCmd.AddCommand(cmdenv.InfoCmd)
	EnvCmd.AddCommand(cmdenv.ListCmd)
	EnvCmd.AddCommand(cmdenv.RunCmd)
	EnvCmd.AddCommand(cmdenv.UpCmd)

}
