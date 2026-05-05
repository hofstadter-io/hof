package cuetils

import (
	"strings"

	"cuelang.org/go/cue"
)

func GetValByEx(ex, pkg string, val cue.Value) cue.Value {
	if ex == "" || ex == "." {
		return val
	} else {
		p := ExToPath(ex, pkg)
		if p.Err() == nil {
			return val.LookupPath(p)
		} else {
			ctx := val.Context()
			return ctx.CompileString(
				ex,
				cue.Filename("--expression:"+ex),
				cue.InferBuiltins(true),
				cue.Scope(val),
			)
		}
	}
}

func ExToPath(ex, pkg string) cue.Path {
	if pkg == "" {
		pkg = "_"
	}
	var sels []cue.Selector
	// assume we can split on dots
	parts := strings.Split(ex, ".")
	for _, part := range parts {
		if strings.HasPrefix(part, "_") {
			sels = append(sels, cue.Hid(part, pkg))
			// fmt.Println("SELS", pkg, sels)
		} else {
			p := cue.ParsePath(part)
			sels = append(sels, p.Selectors()...)
			// fmt.Printf("P: %#+v %v\n", p.Selectors(), p.Err())
		}
	}

	return cue.MakePath(sels...)
}
