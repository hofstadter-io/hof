package runtime

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
)

// RunDir runs the tests in the given directory. All files in dir with a ".hls"
// are considered to be test files.
func Run(t *testing.T, p Params) {
	RunT(tshim{t}, p)
}

// T holds all the methods of the *testing.T type that
// are used by testscript.
type T interface {
	Skip(...interface{})
	Fatal(...interface{})
	Parallel()
	Log(...interface{})
	FailNow()
	Run(string, func(T))
	// Verbose is usually implemented by the testing package
	// directly rather than on the *testing.T type.
	Verbose() bool
}

type tshim struct {
	*testing.T
}

func (t tshim) Run(name string, f func(T)) {
	t.T.Run(name, func(t *testing.T) {
		f(tshim{t})
	})
}

func (t tshim) Verbose() bool {
	return testing.Verbose()
}

// RunT is like Run but uses an interface type instead of the concrete *testing.T
// type to make it possible to use testscript functionality outside of go test.
func RunT(t T, p Params) {
	// add any defaults that were not specified
	p = paramDefaults(p)

	glob := filepath.Join(p.Dir, p.Glob)
	files, err := filepath.Glob(glob)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal(fmt.Sprintf("no scripts found matching glob: %v", glob))
	}

	sort.Strings(files)

	testTempDir := p.WorkdirRoot
	if testTempDir == "" {
		testTempDir, err = ioutil.TempDir(os.Getenv("GOTMPDIR"), "go-test-script")
		if err != nil {
			t.Fatal(err)
		}
	} else {
		p.TestWork = true
	}
	// The temp dir returned by ioutil.TempDir might be a sym linked dir (default
	// behaviour in macOS). That could mess up matching that includes $WORK if,
	// for example, an external program outputs resolved paths. Evaluating the
	// dir here will ensure consistency.
	testTempDir, err = filepath.EvalSymlinks(testTempDir)
	if err != nil {
		t.Fatal(err)
	}
	refCount := int32(len(files))
	for _, file := range files {
		file := file
		name := strings.TrimSuffix(filepath.Base(file), ".hls")
		t.Run(name, func(t T) {
			t.Parallel()
			ts := &Script{
				t:             t,
				testTempDir:   testTempDir,
				name:          name,
				file:          file,
				params:        p,
				ctxt:          context.Background(),
				deferred:      func() {},
				scriptFiles:   make(map[string]string),
				scriptUpdates: make(map[string]string),
			}
			defer func() {
				if ts.failed {
					fmt.Println(ts.log.String())
					os.Exit(1)
				}
				if p.TestWork || *testWork {
					return
				}
				removeAll(ts.workdir)
				if atomic.AddInt32(&refCount, -1) == 0 {
					// This is the last subtest to finish. Remove the
					// parent directory too.
					os.Remove(testTempDir)
				}

			}()
			ts.run()

		})
	}
}

// Defer arranges for f to be called at the end
// of the test. If Defer is called multiple times, the
// defers are executed in reverse order (similar
// to Go's defer statement)
func (ts *Script) Defer(f func()) {
	old := ts.deferred
	ts.deferred = func() {
		defer old()
		f()
	}
}

// Check calls ts.Fatalf if err != nil.
func (ts *Script) Check(err error) {
	if err != nil {
		ts.Fatalf("%v", err)
	}
}

// Logf appends the given formatted message to the test log transcript.
func (ts *Script) Logf(format string, args ...interface{}) {
	format = strings.TrimSuffix(format, "\n")
	fmt.Fprintf(&ts.log, format, args...)
	ts.log.WriteByte('\n')
}

// fatalf aborts the test with the given failure message.
func (ts *Script) Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(&ts.log, "FAIL: %s:%d: %s\n", ts.file, ts.lineno, fmt.Sprintf(format, args...))
	ts.failed = true
	ts.t.FailNow()
}
