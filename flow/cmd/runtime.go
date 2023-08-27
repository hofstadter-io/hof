package cmd

import (
	"fmt"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/flow/flow"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/runtime"
)

// gen.Runtime extends the common runtime.Runtime
type Runtime struct {
	*runtime.Runtime

	// Setup options
	FlowFlags     flags.FlowPflagpole
}

func NewFlowRuntime(RT *runtime.Runtime, cflags flags.FlowPflagpole) (*Runtime) {
	return &Runtime{
		Runtime:  RT,
		FlowFlags: cflags,
	}

}

func prepRuntime(args []string, rflags flags.RootPflagpole, cflags flags.FlowPflagpole) (*Runtime, error) {
	// unsugar the @flow-names into popts
	var entries, flowArgs, tagArgs []string
	for _, e := range args {
		if strings.HasPrefix(e, "@") {
			flowArgs = append(flowArgs, strings.TrimPrefix(e, "@"))
		} else if strings.HasPrefix(e, "+") {
			tagArgs = append(tagArgs, strings.TrimPrefix(e, "+"))
		} else {
			entries = append(entries, e)
		}
	}

	// update entrypoints and Flow flags
	rflags.Tags = append(rflags.Tags, tagArgs...)
	cflags.Flow = append(cflags.Flow, flowArgs...)

	if rflags.Verbosity > 0 {
		fmt.Println("flow modified inputs", entries, cflags.Flow, rflags.Tags)
	}

	// create our core runtime
	r, err := runtime.New(entries, rflags)
	if err != nil {
		return nil, err
	}

	// upgrade to a generator runtime
	R := NewFlowRuntime(r, cflags)

	err = R.Load()
	if err != nil {
		return R, cuetils.ExpandCueError(err)
	}

	if R.Value.Err() != nil {
		fmt.Println("prepRuntime Error:", R.Value.Err())
		return R, cuetils.ExpandCueError(R.Value.Validate())
	}

	err = R.EnrichFlows(cflags.Flow, NoOp)
	if err != nil {
		return R, cuetils.ExpandCueError(err)
	}

	// log cue dirs
	if R.Flags.Verbosity > 1 {
		fmt.Println("CueDirs:", R.CueModuleRoot, R.WorkingDir, R.CwdToRoot)
	}
	if len(R.Workflows) == 0 {
		return R, fmt.Errorf("no workflows found")
	}


	return R, nil
}

func NoOp(R *runtime.Runtime, flow *flow.Flow) error {

	return nil
}

