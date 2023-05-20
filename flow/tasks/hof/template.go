package hof

import (
	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"

	"github.com/hofstadter-io/hof/lib/templates"
)

type HofTemplate struct{
	Name     string
	Data     any
	Template string
	Partials map[string]string

	Delims   templates.Delims
}

func NewHofTemplate(val cue.Value) (hofcontext.Runner, error) {
	return &HofTemplate{}, nil
}

func (T *HofTemplate) Run(ctx *hofcontext.Context) (interface{}, error) {

	v := ctx.Value

	ferr := func() error {
		ctx.CUELock.Lock()
		defer func() {
			ctx.CUELock.Unlock()
		}()

		err := v.Decode(T)
		
		return err
	}()
	if ferr != nil {
		return nil, ferr
	}

	t, err := templates.CreateFromString(T.Name, T.Template, T.Delims)
	if err != nil {
		return nil, err
	}

	for k, P := range T.Partials {
		p := t.T.New(k)
		// do we need to do this, does the partial use the helpers already registered?
		// T.AddGolangHelpers()
		_, err := p.Parse(P)
		if err != nil {
			return nil, err
		}
	}

	bs, err := t.Render(T.Data)
	if err != nil {
		return nil, err
	}


	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()
	res := v.FillPath(cue.ParsePath("out"), string(bs))

	return res, nil
}

