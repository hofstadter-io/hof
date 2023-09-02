package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var RunFlagSet *pflag.FlagSet

type RunFlagpole struct {
	List        bool
	Info        bool
	Suite       []string
	Runner      []string
	Environment []string
	Data        []string
	Workdir     string
}

var RunFlags RunFlagpole

func SetupRunFlags(fset *pflag.FlagSet, fpole *RunFlagpole) {
	// flags

	fset.BoolVarP(&(fpole.List), "list", "", false, "list matching scripts that would run")
	fset.BoolVarP(&(fpole.Info), "info", "", false, "view detailed info for matching scripts")
	fset.StringArrayVarP(&(fpole.Suite), "suite", "", nil, "<name>: _ @run(suite)'s to run")
	fset.StringArrayVarP(&(fpole.Runner), "runner", "", nil, "<name>: _ @run(script)'s to run")
	fset.StringArrayVarP(&(fpole.Environment), "env", "", nil, "exrta environment variables for scripts")
	fset.StringArrayVarP(&(fpole.Data), "data", "", nil, "extra data to include in the scripts context")
	fset.StringVarP(&(fpole.Workdir), "workdir", "w", "", "working directory")
}

func init() {
	RunFlagSet = pflag.NewFlagSet("Run", pflag.ContinueOnError)

	SetupRunFlags(RunFlagSet, &RunFlags)

}
