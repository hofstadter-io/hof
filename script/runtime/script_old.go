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
	goruntime "runtime"
	"strings"

	"github.com/hofstadter-io/hof/lib/gotils/txtar"
)

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

