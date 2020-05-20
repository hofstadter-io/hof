package runtime

import (
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/pflags"
	"github.com/hofstadter-io/hof/gen/cuefig"
	"github.com/hofstadter-io/hof/lib/structural"
)

var rt *Runtime

func init() {
	rt = NewRuntime()
}

func Init() error {
	r := NewRuntime()
	err := r.Init()
	if err != nil {
		return err
	}

	rt = r

	return nil
}

func GetRuntime() *Runtime {
	return rt
}

// Runtime holds the app config/secrets
type Runtime struct {
	ContextType  string
	ContextValue cue.Value
	ConfigType  string
	ConfigValue cue.Value
	SecretType  string
	SecretValue  cue.Value
}

func NewRuntime() *Runtime {
	return &Runtime{}
}
// TODO Load user/app config/secret

// We can safely ignore errors here. If the file exists, cue errors will be printed, otherwise up to the user
func (R *Runtime) Init() (err error) {
	// These are used to track if we found a file or not
	contextFound, configFound, secretFound := false, false, false

	// First check config/secret flags, non-existance should err as user specified a flag
	//  if they exist, we load into local because we prefer that later
	if pflags.RootContextPflag != "" {
		val, err := cuefig.LoadContextConfig("", pflags.RootContextPflag)
		if err != nil {
			// Return early if they specify a file and we don't find it
			return err
		}
		contextFound = true
		R.ContextValue = val
		R.ContextType = "custom-context"
	}
	if pflags.RootConfigPflag != "" {
		val, err := cuefig.LoadConfigConfig("", pflags.RootConfigPflag)
		if err != nil {
			// Return early if they specify a file and we don't find it
			return err
		}
		configFound = true
		R.ConfigValue = val
		R.ConfigType = "custom-config"
	}
	if pflags.RootSecretPflag != "" {
		val, err := cuefig.LoadSecretConfig("", pflags.RootSecretPflag)
		if err != nil {
			// Return early if they specify a file and we don't find it
			return err
		}
		secretFound = true
		R.SecretValue = val
		R.SecretType = "custom-secret"
	}

	// Second, look for local config/secret
	if !contextFound {
		val, err := cuefig.LoadContextDefault()
		// NOTE, we are doing the opposite of normal err checks here
		if err == nil {
			configFound = true
			R.ContextValue = val
			R.ContextType = "local-context"
		}
	}
	if !configFound {
		val, err := cuefig.LoadConfigDefault()
		// NOTE, we are doing the opposite of normal err checks here
		if err == nil {
			configFound = true
			R.ConfigValue = val
			R.ConfigType = "local-config"
		}
	}
	if !secretFound {
		val, err := cuefig.LoadSecretDefault()
		// NOTE, we are doing the opposite of normal err checks here
		if err == nil {
			secretFound = true
			R.SecretValue = val
			R.SecretType = "local-secret"
		}
	}

	// Finally, check for global config/secret
	if !contextFound {
		val, err := cuefig.LoadHofctxDefault()
		// NOTE, we are doing the opposite of normal err checks here
		if err == nil {
			contextFound = true
			R.ContextValue = val
			R.ContextType = "global-context"
		}
	}
	if !configFound {
		val, err := cuefig.LoadHofcfgDefault()
		// NOTE, we are doing the opposite of normal err checks here
		if err == nil {
			configFound = true
			R.ConfigValue = val
			R.ConfigType = "global-config"
		}
	}
	if !secretFound {
		val, err := cuefig.LoadHofshhDefault()
		// NOTE, we are doing the opposite of normal err checks here
		if err == nil {
			secretFound = true
			R.SecretValue = val
			R.SecretType = "global-secret"
		}
	}

	return err
}

func (R *Runtime) PrintConfig() error {
	// Get top level struct from cuelang
	S, err := R.ConfigValue.Struct()
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

func (R *Runtime) PrintSecret() error {
	// Get top level struct from cuelang
	S, err := R.SecretValue.Struct()
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

func (R *Runtime) ContextGet(path string) (cue.Value, error) {
	var orig cue.Value
	var err error
	if pflags.RootContextPflag != "" {
		orig, err = cuefig.LoadContextConfig("", pflags.RootContextPflag)
	} else if pflags.RootLocalPflag {
		orig, err = cuefig.LoadContextConfig("", cuefig.ContextEntrypoint)
	} else if pflags.RootGlobalPflag {
		orig, err = cuefig.LoadHofctxDefault()
	} else {
		orig, err = cuefig.LoadContextDefault()
	}

	// now check for error
	if err != nil {
		return orig, err
	}

	if path == "" {
		return orig, nil
	}
	paths := strings.Split(path, ".")
	val := orig.Lookup(paths...)
	return val, nil
}

func (R *Runtime) ConfigGet(path string) (cue.Value, error) {
	var orig cue.Value
	var err error
	if pflags.RootConfigPflag != "" {
		orig, err = cuefig.LoadConfigConfig("", pflags.RootConfigPflag)
	} else if pflags.RootLocalPflag {
		orig, err = cuefig.LoadConfigConfig("", cuefig.ConfigEntrypoint)
	} else if pflags.RootGlobalPflag {
		orig, err = cuefig.LoadHofcfgDefault()
	} else {
		orig, err = cuefig.LoadConfigDefault()
	}

	// now check for error
	if err != nil {
		return orig, err
	}

	if path == "" {
		return orig, nil
	}
	paths := strings.Split(path, ".")
	val := orig.Lookup(paths...)
	return val, nil
}

func (R *Runtime) SecretGet(path string) (cue.Value, error) {
	var orig cue.Value
	var err error
	if pflags.RootSecretPflag != "" {
		orig, err = cuefig.LoadSecretConfig("", pflags.RootSecretPflag)
	} else if pflags.RootLocalPflag {
		orig, err = cuefig.LoadSecretConfig("", cuefig.SecretEntrypoint)
	} else if pflags.RootGlobalPflag {
		orig, err = cuefig.LoadHofshhDefault()
	} else {
		orig, err = cuefig.LoadSecretDefault()
	}

	// now check for error
	if err != nil {
		return orig, err
	}

	if path == "" {
		return orig, nil
	}
	paths := strings.Split(path, ".")
	val := orig.Lookup(paths...)
	return val, nil
}

func (R *Runtime) ContextSet(expr string) (error) {
	var orig cue.Value
	var val cue.Value
	var err error

	// Check which config we want to work with
	if pflags.RootContextPflag != "" {
		orig, err = cuefig.LoadContextConfig("", pflags.RootContextPflag)
	} else if pflags.RootLocalPflag {
		orig, err = cuefig.LoadContextConfig("", cuefig.ContextEntrypoint)
	} else if pflags.RootGlobalPflag {
		orig, err = cuefig.LoadHofctxDefault()
	} else {
		orig, err = cuefig.LoadContextDefault()
	}

	// now check for error from that config selection process
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
			// error is worse than non-existant
			return err
		}
		// file does not exist, so we should just set
		var r cue.Runtime
		inst, err := r.Compile("", expr)
		if err != nil {
			return err
		}
		val = inst.Value()
		if val.Err() != nil {
			return val.Err()
		}

	} else {
		val, err = structural.Merge(orig, expr)
		if err != nil {
			return err
		}
	}

	// Now save
	if pflags.RootContextPflag != "" {
		err = cuefig.SaveContextConfig("", pflags.RootContextPflag, val)
	} else if pflags.RootLocalPflag {
		err = cuefig.SaveContextConfig("", cuefig.ContextEntrypoint, val)
	} else if pflags.RootGlobalPflag {
		err = cuefig.SaveHofctxDefault(val)
	} else {
		err = cuefig.SaveContextDefault(val)
	}
	return err
}

func (R *Runtime) ConfigSet(expr string) (error) {
	var orig cue.Value
	var val cue.Value
	var err error

	// Check which config we want to work with
	if pflags.RootConfigPflag != "" {
		orig, err = cuefig.LoadConfigConfig("", pflags.RootConfigPflag)
	} else if pflags.RootLocalPflag {
		orig, err = cuefig.LoadConfigConfig("", cuefig.ConfigEntrypoint)
	} else if pflags.RootGlobalPflag {
		orig, err = cuefig.LoadHofcfgDefault()
	} else {
		orig, err = cuefig.LoadConfigDefault()
	}

	// now check for error from that config selection process
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
			// error is worse than non-existant
			return err
		}
		// file does not exist, so we should just set
		var r cue.Runtime
		inst, err := r.Compile("", expr)
		if err != nil {
			return err
		}
		val = inst.Value()
		if val.Err() != nil {
			return val.Err()
		}

	} else {
		val, err = structural.Merge(orig, expr)
		if err != nil {
			return err
		}
	}

	// Now save
	if pflags.RootConfigPflag != "" {
		err = cuefig.SaveConfigConfig("", pflags.RootConfigPflag, val)
	} else if pflags.RootLocalPflag {
		err = cuefig.SaveConfigConfig("", cuefig.ConfigEntrypoint, val)
	} else if pflags.RootGlobalPflag {
		err = cuefig.SaveHofcfgDefault(val)
	} else {
		err = cuefig.SaveConfigDefault(val)
	}
	return err
}

func (R *Runtime) SecretSet(expr string) (error) {
	var orig cue.Value
	var val cue.Value
	var err error

	// Check which config we want to work with
	if pflags.RootSecretPflag != "" {
		orig, err = cuefig.LoadSecretConfig("", pflags.RootSecretPflag)
	} else if pflags.RootLocalPflag {
		orig, err = cuefig.LoadSecretConfig("", cuefig.SecretEntrypoint)
	} else if pflags.RootGlobalPflag {
		orig, err = cuefig.LoadHofshhDefault()
	} else {
		orig, err = cuefig.LoadSecretDefault()
	}

	// now check for error from that config selection process
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
			// error is worse than non-existant
			return err
		}
		// file does not exist, so we should just set
		var r cue.Runtime
		inst, err := r.Compile("", expr)
		if err != nil {
			return err
		}
		val = inst.Value()
		if val.Err() != nil {
			return val.Err()
		}

	} else {
		val, err = structural.Merge(orig, expr)
		if err != nil {
			return err
		}
	}

	// Now save
	if pflags.RootSecretPflag != "" {
		err = cuefig.SaveSecretConfig("", pflags.RootSecretPflag, val)
	} else if pflags.RootLocalPflag {
		err = cuefig.SaveSecretConfig("", cuefig.SecretEntrypoint, val)
	} else if pflags.RootGlobalPflag {
		err = cuefig.SaveHofshhDefault(val)
	} else {
		err = cuefig.SaveSecretDefault(val)
	}
	return err
}
