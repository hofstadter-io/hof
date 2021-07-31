package cuetils

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
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

	CueContext *cue.Context
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
			i := CRT.CueContext.CompileString(T)
			if i.Err() != nil {
				return O, i.Err()
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

	var errs []error

	// XXX TODO XXX
	//  add the second arg from our runtime when implemented
	CRT.CueContext = cuecontext.New()
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
		V := CRT.CueContext.BuildInstance(bi)
		if V.Err() != nil {
			es := errors.Errors(V.Err())
			// fmt.Println("BUILD ERR", es, I)
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}

		CRT.CueValue = V

		// Decode? we want to be lazy
		/*
		err = V.Decode(&CRT.Value)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		*/

		// fmt.Println(i, "decoded", CRT.Value)

	}

	if len(errs) > 0 {
		CRT.CueErrors = errs
		return fmt.Errorf("Errors while loading Cue entrypoints: %s %v\n%v", CRT.Workspace, CRT.Entrypoints, errs)
	}

	return nil
}

