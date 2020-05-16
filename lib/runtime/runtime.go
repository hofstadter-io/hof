package runtime

import (
	"fmt"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/pflags"
	"github.com/hofstadter-io/hof/lib/cuefig"
	"github.com/hofstadter-io/hof/lib/types"
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

type Runtime struct {
	Config       *types.Config
	ConfigCueVal cue.Value
	Creds        *types.Creds
	CredsCueVal  cue.Value
}

func NewRuntime() *Runtime {
	return &Runtime{}
}

func (R *Runtime) Init() error {

	R.Config = &types.Config{}
	var val1, val2 cue.Value
	var err error
	if pflags.RootConfigPflag != "" {
		val1, err = cuefig.LoadConfigConfig(pflags.RootConfigPflag, R.Config)
	} else {
		val1, err = cuefig.LoadConfigDefault(R.Config)
	}
	if pflags.RootCredsPflag != "" {
		val2, err = cuefig.LoadSecretConfig(pflags.RootCredsPflag, R.Creds)
	} else {
		val2, err = cuefig.LoadSecretDefault(R.Creds)
	}

	if err != nil {
		return err
	}
	R.ConfigCueVal = val1
	R.CredsCueVal = val2
	return nil
}

func (R *Runtime) Print() error {
	// Get top level struct from cuelang
	S, err := R.ConfigCueVal.Struct()
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
