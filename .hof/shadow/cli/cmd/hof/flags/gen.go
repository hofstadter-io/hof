package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var GenFlagSet *pflag.FlagSet

type GenFlagpole struct {
	Generator   []string
	Template    []string
	Partial     []string
	Diff3       bool
	NoFormat    bool
	KeepDeleted bool
	Exec        bool
	Watch       bool
	WatchFull   []string
	WatchFast   []string
	AsModule    string
	Outdir      string
}

var GenFlags GenFlagpole

func SetupGenFlags(fset *pflag.FlagSet, fpole *GenFlagpole) {
	// flags

	fset.StringArrayVarP(&(fpole.Generator), "generator", "G", nil, "generator tags to run, default is all, or none if -T is used")
	fset.StringArrayVarP(&(fpole.Template), "template", "T", nil, "template mapping to render, see help for format")
	fset.StringArrayVarP(&(fpole.Partial), "partial", "P", nil, "file globs to partial templates to register with the templates")
	fset.BoolVarP(&(fpole.Diff3), "diff3", "3", false, "enable diff3 support for custom code")
	fset.BoolVarP(&(fpole.NoFormat), "no-format", "", false, "disable formatting during code gen (ad-hoc only)")
	fset.BoolVarP(&(fpole.KeepDeleted), "keep-deleted", "", false, "keep files that would be deleted after code generation")
	fset.BoolVarP(&(fpole.Exec), "exec", "", false, "enable pre/post-exec support when generating code")
	fset.BoolVarP(&(fpole.Watch), "watch", "w", false, "run in watch mode, regenerating when files change, implied by -W/X")
	fset.StringArrayVarP(&(fpole.WatchFull), "watch-globs", "W", nil, "filepath globs to watch for changes and trigger full regen")
	fset.StringArrayVarP(&(fpole.WatchFast), "watch-fast", "X", nil, "filepath globs to watch for changes and trigger fast regen")
	fset.StringVarP(&(fpole.AsModule), "as-module", "", "", "<github.com/username/<name>> like value for the generator module made from the given flags")
	fset.StringVarP(&(fpole.Outdir), "outdir", "O", "", "base directory to write all output u")
}

func init() {
	GenFlagSet = pflag.NewFlagSet("Gen", pflag.ContinueOnError)

	SetupGenFlags(GenFlagSet, &GenFlags)

}
