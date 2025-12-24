package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var EnvFlagSet *pflag.FlagSet

type EnvPflagpole struct {
	Progress    string
	Interactive bool
	NoCache     bool
}

func SetupEnvPflags(fset *pflag.FlagSet, fpole *EnvPflagpole) {
	// pflags

	fset.StringVarP(&(fpole.Progress), "progress", "P", "auto", "output format [auto, plain, tty, dots, report (for ai)]")
	fset.BoolVarP(&(fpole.Interactive), "interactive", "X", false, "enter an interactive terminal on failure, requires a tty")
	fset.BoolVarP(&(fpole.NoCache), "no-cache", "Z", false, "bust the cache and force evaluation")
}

var EnvPflags EnvPflagpole

func init() {
	EnvFlagSet = pflag.NewFlagSet("Env", pflag.ContinueOnError)

	SetupEnvPflags(EnvFlagSet, &EnvPflags)

}
