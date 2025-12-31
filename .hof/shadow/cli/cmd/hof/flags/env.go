package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var EnvFlagSet *pflag.FlagSet

type EnvPflagpole struct {
	Progress  string
	OnFailure bool
	NoExit    bool
	NoCache   bool
	Shell     bool
	Kind      []string
	Sort      []string
}

func SetupEnvPflags(fset *pflag.FlagSet, fpole *EnvPflagpole) {
	// pflags

	fset.StringVarP(&(fpole.Progress), "progress", "P", "auto", "output format [auto, plain, tty, dots, report (for ai)]")
	fset.BoolVarP(&(fpole.OnFailure), "on-failure", "F", false, "on failure, enter an interactive terminal, requires a tty")
	fset.BoolVarP(&(fpole.NoExit), "no-exit", "N", false, "Leave the TUI open after finishing")
	fset.BoolVarP(&(fpole.NoCache), "no-cache", "Z", false, "bust the cache and force evaluation")
	fset.BoolVarP(&(fpole.Shell), "shell", "S", false, "launch a shell after cmd completion for each target")
	fset.StringArrayVarP(&(fpole.Kind), "kind", "k", nil, "kinds to include, defaults to all")
	fset.StringArrayVarP(&(fpole.Sort), "sort", "s", []string{"name"}, "sort columns, default is the order CUE defines")
}

var EnvPflags EnvPflagpole

func init() {
	EnvFlagSet = pflag.NewFlagSet("Env", pflag.ContinueOnError)

	SetupEnvPflags(EnvFlagSet, &EnvPflags)

}
