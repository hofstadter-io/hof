package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var Env__RunFlagSet *pflag.FlagSet

type Env__RunFlagpole struct {
	Command string
}

var Env__RunFlags Env__RunFlagpole

func SetupEnv__RunFlags(fset *pflag.FlagSet, fpole *Env__RunFlagpole) {
	// flags

	fset.StringVarP(&(fpole.Command), "cmd", "c", "", "the command to run, if none by default or to override")
}

func init() {
	Env__RunFlagSet = pflag.NewFlagSet("Env__Run", pflag.ContinueOnError)

	SetupEnv__RunFlags(Env__RunFlagSet, &Env__RunFlags)

}
