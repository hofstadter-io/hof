// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gotooltest_test

import (
	"testing"

	"github.com/hofstadter-io/hof/lib/gotils/gotooltest"
	"github.com/hofstadter-io/hof/lib/gotils/testscript"
)

func TestSimple(t *testing.T) {
	p := testscript.Params{
		Dir: "testdata",
	}
	if err := gotooltest.Setup(&p); err != nil {
		t.Fatal(err)
	}
	testscript.Run(t, p)
}
