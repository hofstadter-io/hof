// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goproxytest_test

import (
	"path/filepath"
	"testing"

	"github.com/hofstadter-io/hof/lib/gotils/goproxytest"
	"github.com/hofstadter-io/hof/lib/gotils/gotooltest"
	"github.com/hofstadter-io/hof/script"
)

func TestScripts(t *testing.T) {
	srv, err := goproxytest.NewServer(filepath.Join("testdata", "mod"), "")
	if err != nil {
		t.Fatalf("cannot start proxy: %v", err)
	}
	p := script.Params{
		Dir: "testdata",
		Setup: func(e *script.Env) error {
			e.Vars = append(e.Vars,
				"GOPROXY="+srv.URL,
				"GONOSUMDB=*",
			)
			return nil
		},
	}
	if err := gotooltest.Setup(&p); err != nil {
		t.Fatal(err)
	}
	script.Run(t, p)
}
