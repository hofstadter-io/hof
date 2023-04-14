package cmd

import (
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func findMaxLabelLen(R *runtime.Runtime, dflags flags.DatamodelPflagpole) int {
	max := 0
	for _, dm := range R.Datamodels {
		m := dm.FindMaxLabelLen(dflags)
		if m > max {
			max = m
		}
	}
	return max
}


