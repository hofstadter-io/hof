package runtime

import (
	"fmt"
	"strings"
)

// Env holds the environment to use at the start of a test script invocation.
type Env struct {
	// WorkDir holds the path to the root directory of the
	// extracted files.
	WorkDir string
	// Vars holds the initial set environment variables that will be passed to the
	// testscript commands.
	Vars []string
	// Cd holds the initial current working directory.
	Cd string
	// Values holds a map of arbitrary values for use by custom
	// testscript commands. This enables Setup to pass arbitrary
	// values (not just strings) through to custom commands.
	Values map[interface{}]interface{}

	ts *Script
}

// Value returns a value from Env.Values, or nil if no
// value was set by Setup.
func (ts *Script) Value(key interface{}) interface{} {
	return ts.values[key]
}

// Defer arranges for f to be called at the end
// of the test. If Defer is called multiple times, the
// defers are executed in reverse order (similar
// to Go's defer statement)
func (e *Env) Defer(f func()) {
	e.ts.Defer(f)
}

// Getenv retrieves the value of the environment variable named by the key. It
// returns the value, which will be empty if the variable is not present.
func (e *Env) Getenv(key string) string {
	key = envvarname(key)
	for i := len(e.Vars) - 1; i >= 0; i-- {
		if pair := strings.SplitN(e.Vars[i], "=", 2); len(pair) == 2 && envvarname(pair[0]) == key {
			return pair[1]
		}
	}
	return ""
}

// Setenv sets the value of the environment variable named by the key. It
// panics if key is invalid.
func (e *Env) Setenv(key, value string) {
	if key == "" || strings.IndexByte(key, '=') != -1 {
		panic(fmt.Errorf("invalid environment variable key %q", key))
	}
	e.Vars = append(e.Vars, key+"="+value)
}

// T returns the t argument passed to the current test by the T.Run method.
// Note that if the tests were started by calling Run,
// the returned value will implement testing.TB.
// Note that, despite that, the underlying value will not be of type
// *testing.T because *testing.T does not implement T.
//
// If Cleanup is called on the returned value, the function will run
// after any functions passed to Env.Defer.
func (e *Env) T() T {
	return e.ts.t
}


