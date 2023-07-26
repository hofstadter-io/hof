package os

import (
	"fmt"
	g_os "os"

	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Getenv struct{}

func NewGetenv(val cue.Value) (hofcontext.Runner, error) {
	return &Getenv{}, nil
}

func (T *Getenv) Run(ctx *hofcontext.Context) (interface{}, error) {
	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()

	v := ctx.Value

	// If a struct, try to fill all fields with matching ENV VAR
	if v.IncompleteKind() == cue.StructKind {
		flds, err := v.Fields()
		if err != nil {
			return nil, err
		}

		for flds.Next() {
			sel := flds.Selector()
			key := sel.String()
			val := g_os.Getenv(key)
			// fmt.Println("pdeps:", t.PathDependencies(cue.MakePath(sel)))

		// HACK, this works around a bug in CUE
		// p := cue.MakePath(sel)
		p := cue.ParsePath(fmt.Sprint(sel))
			v = v.FillPath(p, val)
		}
	} else {
		// otherwise, try to fill a field
		p := v.Path().Selectors()
		sel := p[len(p)-1]
		key := sel.String()
		val := g_os.Getenv(key)
		v = v.FillPath(cue.ParsePath(""), val)
	}

	return v, nil
}
