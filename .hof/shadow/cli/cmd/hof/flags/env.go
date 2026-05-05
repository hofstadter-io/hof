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
	Kind      []string
	Sort      []string
	EnvVar    []string
	EnvFile   []string
	ShhVar    []string
	ShhFile   []string
	EnvAll    bool
	ShowAll   bool
	FailFast  bool
	Unsafe    bool
	Parallel  int
}

func SetupEnvPflags(fset *pflag.FlagSet, fpole *EnvPflagpole) {
	// pflags

	fset.StringVarP(&(fpole.Renderer), "renderer", "R", "auto", "output format [auto, plain, tty, dots, report (for ai)]")
	fset.BoolVarP(&(fpole.OnFailure), "on-failure", "F", false, "on failure, enter an interactive terminal, requires a tty")
	fset.BoolVarP(&(fpole.NoExit), "no-exit", "N", false, "Leave the TUI open after finishing")
	fset.BoolVarP(&(fpole.NoCache), "no-cache", "Z", false, "bust the cache and force evaluation")
	fset.StringArrayVarP(&(fpole.Kind), "kind", "K", nil, "kinds to include, defaults to all")
	fset.StringArrayVarP(&(fpole.Sort), "sort", "S", nil, "sort columns, can be used multiple times")
	fset.StringArrayVarP(&(fpole.EnvVar), "env-var", "", nil, "key=value ENV vars to pass")
	fset.StringArrayVarP(&(fpole.EnvFile), "env-file", "", nil, "path to a file with ENV vars to pass")
	fset.StringArrayVarP(&(fpole.ShhVar), "shh-var", "", nil, "key=value secret ENV vars to pass")
	fset.StringArrayVarP(&(fpole.ShhFile), "shh-file", "", nil, "path to a file with secret ENV vars to pass")
	fset.BoolVarP(&(fpole.EnvAll), "env-all", "", false, "pass os.Env (everything)")
	fset.BoolVarP(&(fpole.ShowAll), "all", "", false, "show all env targets, not just @env() ones")
	fset.BoolVarP(&(fpole.FailFast), "fail-fast", "", false, "fail at first error instead of attempting all targets")
	fset.BoolVarP(&(fpole.Unsafe), "unsafe", "", false, "set insecure root capabilities and privileged nesting, use at your own risk, needed for inception")
	fset.IntVarP(&(fpole.Parallel), "parallel", "P", 1, "number of args or objects to process at once, they may be highly parallel internally")
}

var EnvPflags EnvPflagpole

func init() {
	EnvFlagSet = pflag.NewFlagSet("Env", pflag.ContinueOnError)

	SetupEnvPflags(EnvFlagSet, &EnvPflags)

}
