// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hofstadter-io/hof/script"
)

func TestMain(m *testing.M) {
	os.Exit(script.RunMain(m, map[string]func() int{
		"txtar-x": main1,
	}))
}

func TestScripts(t *testing.T) {
	script.Run(t, script.Params{
		Dir: "testdata",
		Cmds: map[string]func(ts *script.Script, neg bool, args []string){
			"unquote": unquote,
		},
	})
}

func unquote(ts *script.Script, neg bool, args []string) {
	if neg {
		ts.Fatalf("unsupported: ! unquote")
	}
	for _, arg := range args {
		file := ts.MkAbs(arg)
		data, err := ioutil.ReadFile(file)
		ts.Check(err)
		data = bytes.Replace(data, []byte("\n>"), []byte("\n"), -1)
		data = bytes.TrimPrefix(data, []byte(">"))
		err = ioutil.WriteFile(file, data, 0666)
		ts.Check(err)
	}
}
