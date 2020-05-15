package structural

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
)

type CueRuntime struct {
	Entrypoints []string

	CueRT *cue.Runtime
	BIS   []*build.Instance

	FieldOpts []cue.Option
	CueInstance    *cue.Instance
	TopLevelValue  cue.Value
	TopLevelStruct *cue.Struct
}

func NewCueRuntime() *CueRuntime {
	return &CueRuntime{
		FieldOpts: []cue.Option{
			cue.Attributes(true),
			cue.Concrete(false),
			cue.Definitions(true),
			cue.Docs(true),
			cue.Hidden(true),
			cue.Optional(true),
		},
	}


}

func (CR *CueRuntime) LoadCue(entrypoints []string) []error {
	var errs []error

	CR.Entrypoints = entrypoints

	CR.CueRT = &cue.Runtime{}
	CR.BIS = load.Instances(CR.Entrypoints, &load.Config{
		Package: "",
		Tools: false,
	})

	// fmt.Println("len(BIS):", len(CR.BIS))

	// Build the Instances
	I, err := CR.CueRT.Build(CR.BIS[0])
	if err != nil {
		es := errors.Errors(err)
		// fmt.Println("BUILD ERR", es, I)
		for _, e := range es {
			errs = append(errs, e.(error))
		}
		return errs
	}

	CR.CueInstance = I

	// Get top level value from cuelang
	V := I.Value()
	CR.TopLevelValue = V

	// Get top level struct from cuelang
	S, err := V.Struct()
	if err != nil {
		// fmt.Println("STRUCT ERR", err)
		es := errors.Errors(err)
		for _, e := range es {
			errs = append(errs, e.(error))
		}
		return errs
	}

	CR.TopLevelStruct = S

	return errs
}

func (CR *CueRuntime) PrintValue() error {
	node := CR.TopLevelValue.Syntax(CR.FieldOpts...)

	bytes, err := format.Node(
		node,
		format.TabIndent(false),
		format.UseSpaces(2),
	)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	return nil
}

func PrintCueValue(val cue.Value) error {
	node := val.Syntax(
		cue.Attributes(true),
		cue.Concrete(false),
		cue.Definitions(true),
		cue.Docs(true),
		cue.Hidden(true),
		cue.Optional(true),
	)

	bytes, err := format.Node(
		node,
		format.TabIndent(false),
		format.UseSpaces(2),
	)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	return nil
}

