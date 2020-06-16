// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Script-driven tests.
// See testdata/script/README for an overview.

package runtime

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	goruntime "runtime"
	"strings"
	"testing"

	"github.com/hofstadter-io/hof/lib/gotils/imports"
	"github.com/hofstadter-io/hof/lib/gotils/intern/os/execpath"
	"github.com/hofstadter-io/hof/lib/gotils/par"
	"github.com/hofstadter-io/hof/lib/gotils/testenv"
	"github.com/hofstadter-io/hof/lib/gotils/txtar"
)

var execCache par.Cache

// If -testwork is specified, the test prints the name of the temp directory
// and does not remove it when done, so that a programmer can
// poke at the test file tree afterward.
var testWork = flag.Bool("testwork", false, "")


func (ts *Script) applyScriptUpdates() {
	if len(ts.scriptUpdates) == 0 {
		return
	}
	for name, content := range ts.scriptUpdates {
		found := false
		for i := range ts.archive.Files {
			f := &ts.archive.Files[i]
			if f.Name != name {
				continue
			}
			data := []byte(content)
			if txtar.NeedsQuote(data) {
				data1, err := txtar.Quote(data)
				if err != nil {
					ts.t.Fatal(fmt.Sprintf("cannot update script file %q: %v", f.Name, err))
					continue
				}
				data = data1
			}
			f.Data = data
			found = true
		}
		// Sanity check.
		if !found {
			panic("script update file not found")
		}
	}
	if err := ioutil.WriteFile(ts.file, txtar.Format(ts.archive), 0666); err != nil {
		ts.t.Fatal("cannot update script: ", err)
	}
	ts.Logf("%s updated", ts.file)
}

// condition reports whether the given condition is satisfied.
func (ts *Script) condition(cond string) (bool, error) {
	switch cond {
	case "short":
		return testing.Short(), nil
	case "net":
		return testenv.HasExternalNetwork(), nil
	case "link":
		return testenv.HasLink(), nil
	case "symlink":
		return testenv.HasSymlink(), nil
	case goruntime.GOOS, goruntime.GOARCH:
		return true, nil
	default:
		if imports.KnownArch[cond] || imports.KnownOS[cond] {
			return false, nil
		}
		if strings.HasPrefix(cond, "exec:") {
			prog := cond[len("exec:"):]
			ok := execCache.Do(prog, func() interface{} {
				_, err := execpath.Look(prog, ts.Getenv)
				return err == nil
			}).(bool)
			return ok, nil
		}
		if ts.params.Condition != nil {
			return ts.params.Condition(cond)
		}
		ts.Fatalf("unknown condition %q", cond)
		panic("unreachable")
	}
}

// Helpers for command implementations.

// abbrev abbreviates the actual work directory in the string s to the literal string "$WORK".
func (ts *Script) abbrev(s string) string {
	if ts.params.Mode != "test" {
		return s
	}
	s = strings.Replace(s, ts.workdir, "$WORK", -1)
	if *testWork || ts.params.TestWork {
		// Expose actual $WORK value in environment dump on first line of work script,
		// so that the user can find out what directory -testwork left behind.
		s = "WORK=" + ts.workdir + "\n" + strings.TrimPrefix(s, "WORK=$WORK\n")
	}
	return s
}

// expand applies environment variable expansion to the string s.
func (ts *Script) expand(s string) string {
	return os.Expand(s, func(key string) string {
		if key1 := strings.TrimSuffix(key, "@R"); len(key1) != len(key) {
			return regexp.QuoteMeta(ts.Getenv(key1))
		}
		return ts.Getenv(key)
	})
}

// MkAbs interprets file relative to the test script's current directory
// and returns the corresponding absolute path.
func (ts *Script) MkAbs(file string) string {
	if filepath.IsAbs(file) {
		return file
	}
	return filepath.Join(ts.cd, file)
}

// ReadFile returns the contents of the file with the
// given name, intepreted relative to the test script's
// current directory. It interprets "stdout" and "stderr" to
// mean the standard output or standard error from
// the most recent exec or wait command respectively.
//
// If the file cannot be read, the script fails.
func (ts *Script) ReadFile(file string) string {
	switch file {
	case "stdout":
		return ts.stdout
	case "stderr":
		return ts.stderr
	default:
		file = ts.MkAbs(file)
		data, err := ioutil.ReadFile(file)
		ts.Check(err)
		return string(data)
	}
}

// Setenv sets the value of the environment variable named by the key.
func (ts *Script) Setenv(key, value string) {
	ts.env = append(ts.env, key+"="+value)
	ts.envMap[envvarname(key)] = value
}

// Getenv gets the value of the environment variable named by the key.
func (ts *Script) Getenv(key string) string {
	return ts.envMap[envvarname(key)]
}

func removeAll(dir string) error {
	// module cache has 0444 directories;
	// make them writable in order to remove content.
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // ignore errors walking in file system
		}
		if info.IsDir() {
			os.Chmod(path, 0777)
		}
		return nil
	})
	return os.RemoveAll(dir)
}

func homeEnvName() string {
	switch goruntime.GOOS {
	case "windows":
		return "USERPROFILE"
	case "plan9":
		return "home"
	default:
		return "HOME"
	}
}

func tempEnvName() string {
	switch goruntime.GOOS {
	case "windows":
		return "TMP"
	case "plan9":
		return "TMPDIR" // actually plan 9 doesn't have one at all but this is fine
	default:
		return "TMPDIR"
	}
}

