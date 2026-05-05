package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var Agent__ChatFlagSet *pflag.FlagSet

type Agent__ChatPflagpole struct {
	Sid string
}

func SetupAgent__ChatPflags(fset *pflag.FlagSet, fpole *Agent__ChatPflagpole) {
	// pflags

	fset.StringVarP(&(fpole.Sid), "sid", "", "", "session id to continue from")
}

var Agent__ChatPflags Agent__ChatPflagpole

func init() {
	Agent__ChatFlagSet = pflag.NewFlagSet("Agent__Chat", pflag.ContinueOnError)

	SetupAgent__ChatPflags(Agent__ChatFlagSet, &Agent__ChatPflags)

}
