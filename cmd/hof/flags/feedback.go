package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var FeedbackFlagSet *pflag.FlagSet

type FeedbackPflagpole struct {
	Issue  bool
	Labels string
}

func SetupFeedbackPflags(fset *pflag.FlagSet, fpole *FeedbackPflagpole) {
	// pflags

	fset.BoolVarP(&(fpole.Issue), "issue", "", false, "create an issue (discussion is default)")
	fset.StringVarP(&(fpole.Labels), "labels", "", "feedback", "labels,comma,separated")
}

var FeedbackPflags FeedbackPflagpole

func init() {
	FeedbackFlagSet = pflag.NewFlagSet("Feedback", pflag.ContinueOnError)

	SetupFeedbackPflags(FeedbackFlagSet, &FeedbackPflags)

}
