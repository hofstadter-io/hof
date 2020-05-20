package runtime

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"

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
	configFound, secretFound := false, false

	// First check config/secret flags, non-existance should err as user specified a flag
	//  if they exist, we load into local because we prefer that later
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

func (R *Runtime) ConfigGet(path string) (cue.Value, error) {
	var orig cue.Value
	var err error
	if pflags.RootConfigPflag != "" {
		orig, err = cuefig.LoadConfigConfig("", pflags.RootConfigPflag)
	} else if pflags.RootLocalPflag {
		orig, err = cuefig.LoadConfigDefault()
	} else if pflags.RootGlobalPflag {
		orig, err = cuefig.LoadHofcfgDefault()
	} else {
		orig = R.ConfigValue
	}

	// now check for error
	if err != nil {
		return orig, err
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
		orig, err = cuefig.LoadSecretDefault()
	} else if pflags.RootGlobalPflag {
		orig, err = cuefig.LoadHofshhDefault()
	} else {
		orig = R.SecretValue
	}

	// now check for error
	if err != nil {
		return orig, err
	}

	paths := strings.Split(path, ".")
	val := orig.Lookup(paths...)
	return val, nil
}

func (R *Runtime) ConfigSet(expr string) (error) {
	var orig cue.Value
	var err error
	if pflags.RootConfigPflag != "" {
		orig, err = cuefig.LoadConfigConfig("", pflags.RootConfigPflag)
	} else if pflags.RootLocalPflag {
		orig, err = cuefig.LoadConfigDefault()
	} else if pflags.RootGlobalPflag {
		orig, err = cuefig.LoadHofcfgDefault()
	} else {
		orig = R.ConfigValue
	}

	// now check for error
	if err != nil {
		return err
	}

	val, err := structural.Merge(orig, expr)
	if err != nil {
		return err
	}

	bytes, err := format.Node(val.Syntax())
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))

	// Now save
	/*
	if pflags.RootConfigPflag != "" {
		orig, err = cuefig.LoadConfigConfig("", pflags.RootConfigPflag)
	} else if pflags.RootLocalPflag {
		orig, err = cuefig.LoadConfigDefault()
	} else if pflags.RootGlobalPflag {
		orig, err = cuefig.LoadHofcfgDefault()
	} else {
		orig = R.ConfigValue
	}
	*/
	return nil
}

func (R *Runtime) SecretSet(path string, str string) (error) {

	return nil
}
