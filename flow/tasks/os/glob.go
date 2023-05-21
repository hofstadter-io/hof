package os

import (
	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/lib/yagu"
)

type Glob struct{}

func NewGlob(val cue.Value) (hofcontext.Runner, error) {
	return &Glob{}, nil
}

func (T *Glob) Run(ctx *hofcontext.Context) (interface{}, error) {

	val := ctx.Value

	patterns, err := extractGlobConfig(ctx, val)
	if err != nil {
		return nil, err
	}

	filepaths, err := yagu.FilesFromGlobs(patterns)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"filepaths": filepaths}, nil
}

func extractGlobConfig(ctx *hofcontext.Context, val cue.Value) (patterns []string, err error) {
	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()

	ps := val.LookupPath(cue.ParsePath("globs"))
	if ps.Err() != nil {
		return nil, ps.Err()
	}

	iter, err := ps.List()
	if err != nil {
		return nil, err
	}

	for iter.Next() {
		gv := iter.Value()
		if gv.Err() != nil {
			return nil, gv.Err()
		}
		gs, err := gv.String()
		if err != nil {
			return nil, err
		}

		patterns = append(patterns, gs)
	}

	return patterns, nil
}
