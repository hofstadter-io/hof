package runtime

import (
	"fmt"
	"strings"
	"time"

	"github.com/hofstadter-io/hof/script/ast"
)

func (RT *Runtime) RunPhase(ph *ast.Phase, parent *ast.Result) (r *ast.Result, err error) {
	RT.phase = ph

	// Prep result
	r = ast.NewResult(ph, parent)
	RT.SetMultiWriters(r)

	// TODO convert to RT.logger
	// fmt.Fprintln(RT.Stdout, strings.Repeat("%", ph.Level()), "Phase:", ph.Level(), ph.Title(), ph.DocLine(), ph.BegLine(), ph.EndLine())
	fmt.Fprintln(RT.Stdout, strings.Repeat("%", ph.Level()), ph.Title())

	// start result
	r.BegTime = time.Now()
	defer func() {
		if r.EndTime.IsZero() {
			r.EndTime = time.Now()
		}
	}()

	for _, node := range ph.Nodes() {
		nr, err := RT.RunNode(node, r)
		r.AddResult(nr)
		if err != nil {
			r.AddError(err)
			if RT.params.Mode == Run {
				break
			}
		}
		if RT.stopped {
			break
		}
	}

	r.EndTime = time.Now()
	if len(r.Errors) == 0 {
		r.Status = 0
	} else {
		for _, e := range r.Errors {
			RT.logger.Error(e)
		}
		err = fmt.Errorf("%d Phase errors occurred", len(r.Errors))
		r.Status = 1
	}

	// TODO print none/some/all/etc... based on config
	fmt.Fprintf(RT.Stdout, "[phase:%d] [status:%d] [time:%v]\n", ph.Level(), r.Status, r.EndTime.Sub(r.BegTime))

	return r, err
}


