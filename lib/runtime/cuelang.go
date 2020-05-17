package runtime

// Name: config

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

type CueRuntime struct {
	Entrypoints []string
	Workspace   string
	FS billy.Filesystem

	CueRuntime *cue.Runtime
	BuildInstances []*build.Instance
	CueErrors []error

	CueInstance *cue.Instance
	CueValue cue.Value
	Value    interface{}

}

// CueRuntimeFromArgs builds up a CueRuntime
//  by processing the args passed in
func CueRuntimeFromArgs(args []string) (crt *CueRuntime, err error) {
	crt = &CueRuntime{
		Entrypoints: args,
	}

	err = crt.Load()

	return crt, err
}

// CueRuntimeFromArgsAndFlags builds up a CueRuntime
//  by processing the args passed in AND the current flag values
func CueRuntimeFromArgsAndFlags(args []string) (crt *CueRuntime, err error) {
	crt = &CueRuntime{
		Entrypoints: args,
	}

	// XXX TODO XXX
	// Buildup out arg to load.Instances second arg
	// Add this configuration to our runtime struct

	err = crt.Load()

	return crt, err
}

func (CRT *CueRuntime) Load() (err error) {
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

	fmt.Println("CRT.Load", CRT.Workspace, CRT.Entrypoints)

	var errs []error

	// XXX TODO XXX
	//  add the second arg from our runtime when implemented
	CRT.CueRuntime = &cue.Runtime{}
	CRT.BuildInstances = load.Instances(CRT.Entrypoints, nil)
	for i, bi := range CRT.BuildInstances {
		fmt.Println("%d: start", i)

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

		fmt.Println(i, "built", I)

		CRT.CueInstance = I

		// Get top level value from cuelang
		V := I.Value()
		CRT.CueValue = V
		fmt.Println(i, "valued", V)

		// Decode? we want to be lazy
		err = V.Decode(&CRT.Value)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		fmt.Println(i, "decoded", CRT.Value)

	}

	if len(errs) > 0 {
		CRT.CueErrors = errs
		return fmt.Errorf("Errors while loading: %s %v", CRT.Workspace, CRT.Entrypoints)
	}

	return nil
}

func (CRT *CueRuntime) PrintValue() error {
	// Get top level struct from cuelang
	S, err := CRT.CueValue.Struct()
	if err != nil {
		return err
	}

	iter := S.Fields()
	for iter.Next() {

		label := iter.Label()
		value := iter.Value()
		fmt.Println("  -", label, value)
		for attrKey, attrVal := range value.Attributes() {
			fmt.Println("  --", attrKey)
			for i := 0; i < 5; i++ {
				str, err := attrVal.String(i)
				if err != nil {
					break
				}
				fmt.Println("  ---", str)
			}
		}
	}

	return nil
}
