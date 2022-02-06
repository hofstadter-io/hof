package structural

import (
	"fmt"

	"cuelang.org/go/cue"
)

type Options struct {
  AllTypeErrors  bool
  NodeTypeErrors bool
}

func newAny(ctx *cue.Context) cue.Value {
	return ctx.CompileString("_")
}

func newStruct(ctx *cue.Context) cue.Value {
	return ctx.CompileString("{...}")
}

type listProc int

const (
	LIST_ERR listProc = iota // unknown arg
	LIST_AND                 // all apply
	LIST_OR                  // any apply
	LIST_PER                 // pairwise
)

func getListProcType(val cue.Value) (listProc, error) {
	attr := val.Attribute("list")
	// no attribute or no arg, so default is LIST_AND
	if attr.Err() != nil || attr.NumArgs() == 0 {
		return LIST_AND, nil
	}
	a, err := attr.String(0)
	if err != nil {
		return LIST_ERR, err
	}
	// otherwise, check what arg we have
	switch a {
	case "and", "":
		return LIST_AND, nil
	case "or":
		return LIST_OR, nil
	case "pairwise", "per":
		return LIST_PER, nil
	default:
		return LIST_ERR, fmt.Errorf("Unknown list processing type %q at %v", a, val.Pos())
	}
}

func GetLabel(val cue.Value) cue.Selector {
	ss := val.Path().Selectors()
	s := ss[len(ss)-1]
	return s
}

