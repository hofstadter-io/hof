// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hofstadter-io/hof/lib/gotils/gotooltest"
	"github.com/hofstadter-io/hof/lib/gotils/intern/os/execpath"
	"github.com/hofstadter-io/hof/script"
)

func TestMain(m *testing.M) {
	os.Exit(script.RunMain(m, map[string]func() int{
		"testscript": main1,
	}))
}

func TestScripts(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Fatalf("need go in PATH for these tests")
	}

	var stderr bytes.Buffer
	cmd := exec.Command("go", "env", "GOMOD")
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to run %v: %v\n%s", strings.Join(cmd.Args, " "), err, stderr.String())
	}
	gomod := string(out)

	if gomod == "" {
		t.Fatalf("apparently we are not running in module mode?")
	}

	p := script.Params{
		Dir: "testdata",
		Setup: func(env *script.Env) error {
			env.Vars = append(env.Vars,
				"GOINTERNALMODPATH="+filepath.Dir(gomod),
				"GONOSUMDB=*",
			)
			return nil
		},
		Cmds: map[string]func(ts *script.Script, neg bool, args []string){
			"dropgofrompath": dropgofrompath,
			"setfilegoproxy": setfilegoproxy,
		},
	}
	if err := gotooltest.Setup(&p); err != nil {
		t.Fatal(err)
	}
	script.Run(t, p)
}

func dropgofrompath(ts *script.Script, neg bool, args []string) {
	if neg {
		ts.Fatalf("unsupported: ! dropgofrompath")
	}
	var newPath []string
	for _, d := range filepath.SplitList(ts.Getenv("PATH")) {
		getenv := func(k string) string {
			if k == "PATH" {
				return d
			}
			return ts.Getenv(k)
		}
		if _, err := execpath.Look("go", getenv); err != nil {
			newPath = append(newPath, d)
		}
	}
	ts.Setenv("PATH", strings.Join(newPath, string(filepath.ListSeparator)))
}

func setfilegoproxy(ts *script.Script, neg bool, args []string) {
	if neg {
		ts.Fatalf("unsupported: ! setfilegoproxy")
	}
	path := args[0]
	path = filepath.ToSlash(path)
	// probably sufficient to just handle spaces
	path = strings.Replace(path, " ", "%20", -1)
	if runtime.GOOS == "windows" {
		path = "/" + path
	}
	ts.Setenv("GOPROXY", "file://"+path)
}
