// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"

	"github.com/hofstadter-io/hof/script"
)

func TestMain(m *testing.M) {
	os.Exit(script.RunMain(m, map[string]func() int{
		"txtar-c": main1,
	}))
}

func TestScripts(t *testing.T) {
	script.Run(t, script.Params{
		Dir: "testdata",
	})
}
