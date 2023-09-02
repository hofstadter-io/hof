package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var CreateFlagSet *pflag.FlagSet

type CreateFlagpole struct {
	Generator []string
	Outdir    string
	Exec      bool
}

var CreateFlags CreateFlagpole

func SetupCreateFlags(fset *pflag.FlagSet, fpole *CreateFlagpole) {
	// flags

	fset.StringArrayVarP(&(fpole.Generator), "generator", "G", nil, "generator tags to run, default is all")
	fset.StringVarP(&(fpole.Outdir), "outdir", "O", "", "base directory to write all output to")
	fset.BoolVarP(&(fpole.Exec), "exec", "", false, "enable pre/post-exec support when generating code")
}

func init() {
	CreateFlagSet = pflag.NewFlagSet("Create", pflag.ContinueOnError)

	SetupCreateFlags(CreateFlagSet, &CreateFlags)

}
