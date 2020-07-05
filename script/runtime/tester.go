package runtime

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hofstadter-io/hof/script/ast"
)

// If -testwork is specified, the test prints the name of the temp directory
// and does not remove it when done, so that a programmer can
// poke at the test file tree afterward.
var testWork = flag.Bool("testwork", false, "")

// RunDir runs the tests in the given directory. All files in dir with a ".hls"
// are considered to be test files.
func RunTester(t *testing.T, p Params) {
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

			t.Log(name)

			parser := ast.NewParser(nil)

			params := &Params{
				Mode: Test,
			}
			RT := &Runtime{
				t: t,
				params: params,
				parser: parser,
				ctxt:          context.Background(),
				deferred:      func() {},
				stdinR: os.Stdin,
				stdout: os.Stdout,
				stderr: os.Stderr,
				currdir: testTempDir,
				workdir: testTempDir,
				// TODO, clean this up?
				logger: parser.GetLogger(),
			}

			// RT.setupLogger()

			S, err := parser.ParseScript(name)
			if err != nil {
				RT.logger.Error(err)
				t.Fatal(err)
			}

			RT.script = S
			RT.setupEnv()

			defer func() {
				if p.TestWork || *testWork {
					return
				}
				removeAll(RT.workdir)
				if atomic.AddInt32(&refCount, -1) == 0 {
					// This is the last subtest to finish. Remove the
					// parent directory too.
					os.Remove(testTempDir)
				}
			}()

			err = RT.Run()
			if err != nil {
				RT.logger.Error(err)
				t.Fatal(err)
			}

			if RT.failed {
				t.Fatal(err)
				// ^^^ this was not here
				// Was this just for failing out of run mode in the old scripts?
				// os.Exit(1)
			}

			// TODO, check for errors / "exit != 0""
			// perhaps only when in run mode?
			// (need to think about this, likely configurable with sensible defaults for test v run v shell)
		})// see also cmd_exec.go:58
	}
}

// Defer arranges for f to be called at the end
// of the test. If Defer is called multiple times, the
// defers are executed in reverse order (similar
// to Go's defer statement)
func (RT *Runtime) Defer(f func()) {
	old := RT.deferred
	RT.deferred = func() {
		defer old()
		f()
	}
}

// Check calls ts.Fatalf if err != nil.
func (RT *Runtime) Check(err error) {
	if err != nil {
		RT.Fatalf("%v", err)
	}
}

// Logf appends the given formatted message to the test log transcript.
func (RT *Runtime) Logf(format string, args ...interface{}) {
	format = strings.TrimSuffix(format, "\n")
	fmt.Fprintf(RT.stdout, format, args...)
}

// fatalf aborts the test with the given failure message.
func (RT *Runtime) Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(RT.stdout, "FAIL: %s\n", fmt.Sprintf(format, args...))

	if RT.params.Mode == Run {
		RT.stopped = true
		RT.failed = true
	}

	if RT.params.Mode == Test {
		RT.t.FailNow()
	}
}

func paramDefaults(p Params) Params {
	if p.Mode == None {
		p.Mode = Test
	}

	if p.Glob == "" {
		p.Glob = "*.hls"
	}

	return p
}

