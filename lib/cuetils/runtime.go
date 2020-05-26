package cuetils

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
)

type CueSyntaxOptions struct {
	Attributes  bool
	Concrete    bool
	Definitions bool
	Docs        bool
	Hidden      bool
	Optional    bool
}

func (CSO CueSyntaxOptions) MakeOpts() []cue.Option {
	return []cue.Option{
		cue.Attributes(CSO.Attributes),
		cue.Concrete(CSO.Concrete),
		cue.Definitions(CSO.Definitions),
		cue.Docs(CSO.Docs),
		cue.Hidden(CSO.Hidden),
		cue.Optional(CSO.Optional),
	}
}

var (
	DefaultSyntaxOpts = CueSyntaxOptions{
		Attributes: true,
		Concrete: false,
		Definitions: true,
		Docs: true,
		Hidden: true,
		Optional: true,
	}
)

type CueRuntime struct {

	Entrypoints []string
	Workspace   string
	FS billy.Filesystem

	CueRuntime *cue.Runtime
	CueConfig *load.Config
	BuildInstances []*build.Instance
	CueErrors []error
	FieldOpts []cue.Option

	CueInstance *cue.Instance
	CueValue cue.Value
	Value    interface{}

}

func (CRT *CueRuntime) ConvertToValue(in interface{}) (cue.Value, error) {
	O, ook := in.(cue.Value)
	if !ook {
		switch T := in.(type) {
		case string:
			i, err := CRT.CueRuntime.Compile("", in)
			if err != nil {
				return O, err
			}
			v := i.Value()
			if v.Err() != nil {
				return v, v.Err()
			}
			O = v

		default:
			return O, fmt.Errorf("unknown type %v in convertToValue(in)", T)
		}
	}

	return O, nil
}

func (CRT *CueRuntime) Load() (err error) {
	return CRT.load()
}

func (CRT *CueRuntime) load() (err error) {
	// possibly, check for workpath
	if CRT.Workspace != "" {
		_, err = os.Lstat(CRT.Workspace)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
				// error is worse than non-existant
				return err
			}
			// otherwise, does not exist, so we should init?
			// XXX want to let applications decide how to handle this
			return err
		}
	}

	// fmt.Println("CRT.Load", CRT.Workspace, CRT.Entrypoints)

	var errs []error

	// XXX TODO XXX
	//  add the second arg from our runtime when implemented
	CRT.CueRuntime = &cue.Runtime{}
	CRT.BuildInstances = load.Instances(CRT.Entrypoints, nil)
	for _, bi := range CRT.BuildInstances {
		// fmt.Printf("%d: start\n", i)

		if bi.Err != nil {
			fmt.Println("BI ERR", bi.Err, bi.Incomplete, bi.DepsErrors)
			es := errors.Errors(bi.Err)
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}

		// Build the Instance
		I, err := CRT.CueRuntime.Build(bi)
		if err != nil {
			es := errors.Errors(err)
			// fmt.Println("BUILD ERR", es, I)
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}

		// fmt.Println(i, "built", I)

		CRT.CueInstance = I

		// Get top level value from cuelang
		V := I.Value()
		CRT.CueValue = V
		// fmt.Println(i, "valued", V)

		// Decode? we want to be lazy
		err = V.Decode(&CRT.Value)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// fmt.Println(i, "decoded", CRT.Value)

	}

	if len(errs) > 0 {
		CRT.CueErrors = errs
		return fmt.Errorf("Errors while loading: %s %v", CRT.Workspace, CRT.Entrypoints)
	}

	return nil
}

