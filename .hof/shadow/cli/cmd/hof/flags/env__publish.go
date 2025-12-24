package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var Env__PublishFlagSet *pflag.FlagSet

type Env__PublishFlagpole struct {
	Registry string
	Tag      []string
}

var Env__PublishFlags Env__PublishFlagpole

func SetupEnv__PublishFlags(fset *pflag.FlagSet, fpole *Env__PublishFlagpole) {
	// flags

	fset.StringVarP(&(fpole.Registry), "registry", "R", "host.docker.internal:5000", "registry to push to, defaults to veg internal")
	fset.StringArrayVarP(&(fpole.Tag), "tag", "T", []string{"local"}, "tags to give to the environment, can be set multiple times")
}

func init() {
	Env__PublishFlagSet = pflag.NewFlagSet("Env__Publish", pflag.ContinueOnError)

	SetupEnv__PublishFlags(Env__PublishFlagSet, &Env__PublishFlags)

}
