package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var EnvFlagSet *pflag.FlagSet

type EnvPflagpole struct {
	Renderer  string
	OnFailure bool
	NoExit    bool
	NoCache   bool
	Path      []string
	Kind      []string
	Sort      []string
}

func SetupEnvPflags(fset *pflag.FlagSet, fpole *EnvPflagpole) {
	// pflags

	fset.StringVarP(&(fpole.Renderer), "renderer", "R", "auto", "output format [auto, plain, tty, dots, report (for ai)]")
	fset.BoolVarP(&(fpole.OnFailure), "on-failure", "F", false, "on failure, enter an interactive terminal, requires a tty")
	fset.BoolVarP(&(fpole.NoExit), "no-exit", "N", false, "Leave the TUI open after finishing")
	fset.BoolVarP(&(fpole.NoCache), "no-cache", "Z", false, "bust the cache and force evaluation")
	fset.StringArrayVarP(&(fpole.Path), "path", "P", nil, "(cue) path prefixes to include, defaults to all")
	fset.StringArrayVarP(&(fpole.Kind), "kind", "K", nil, "kinds to include, defaults to all")
	fset.StringArrayVarP(&(fpole.Sort), "sort", "S", []string{"name"}, "sort columns, can be used multiple times")
}

var EnvPflags EnvPflagpole

func init() {
	EnvFlagSet = pflag.NewFlagSet("Env", pflag.ContinueOnError)

	SetupEnvPflags(EnvFlagSet, &EnvPflags)

}
