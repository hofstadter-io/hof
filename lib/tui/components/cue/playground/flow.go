package playground

import (
	"bytes"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	flowcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/middleware"
	"github.com/hofstadter-io/hof/flow/tasks"
	"github.com/hofstadter-io/hof/flow/flow"
)


func (C *Playground) runFlow(val cue.Value) (cue.Value, error) {
	var stdin, stdout, stderr bytes.Buffer

	ctx := flowcontext.New()
	ctx.RootValue = val
	ctx.Stdin = &stdin
	ctx.Stdout = &stdout
	ctx.Stderr = &stderr

	// how to inject tags into original value
	// fill / return value
	middleware.UseDefaults(ctx, flags.RootPflagpole{}, flags.FlowPflagpole{})
	tasks.RegisterDefaults(ctx)

	f, err := flow.OldFlow(ctx, val)
	if err != nil {
		return val, err
	}

	err = f.Start()

	return f.Final, err
}
