// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.18

package runtime

import (
	"flag"
	"log"
	"os"
	"testing"
)

// runCoverSubcommand runs the given function, then writes any generated
// coverage information to the cprof file.
// This is called inside a separately run executable.
func runCoverSubcommand(cprof string, mainf func() int) (exitCode int) {
	// Change the error handling mode to PanicOnError
	// so that in the common case of calling flag.Parse in main we'll
	// be able to catch the panic instead of just exiting.
	flag.CommandLine.Init(flag.CommandLine.Name(), flag.PanicOnError)
	defer func() {
		panicErr := recover()
		if _, ok := panicErr.(error); ok {
			// The flag package will already have printed this error, assuming,
			// that is, that the error was created in the flag package.
			// TODO check the stack to be sure it was actually raised by the flag package.
			exitCode = 2
			panicErr = nil
		}
		// Set os.Args so that flag.Parse will tell testing the correct
		// coverprofile setting. Unfortunately this isn't sufficient because
		// the testing oackage explicitly avoids calling flag.Parse again
		// if flag.Parsed returns true, so we the coverprofile value directly
		// too.
		os.Args = []string{os.Args[0], "-test.coverprofile=" + cprof}
		setCoverProfile(cprof)

		// Suppress the chatty coverage and test report.
		devNull, err := os.Open(os.DevNull)
		if err != nil {
			panic(err)
		}
		os.Stdout = devNull
		os.Stderr = devNull

		// Run MainStart (recursively, but it we should be ok) with no tests
		// so that it writes the coverage profile.
		// go1.18 -- m := testing.MainStart(nopTestDeps{}, nil, nil, nil, nil)
		m := testing.MainStart(nopTestDeps{}, nil, nil, nil)
		if code := m.Run(); code != 0 && exitCode == 0 {
			exitCode = code
		}
		if _, err := os.Stat(cprof); err != nil {
			log.Printf("failed to write coverage profile %q", cprof)
		}
		if panicErr != nil {
			// The error didn't originate from the flag package (we know that
			// flag.PanicOnError causes an error value that implements error),
			// so carry on panicking.
			panic(panicErr)
		}
	}()
	return mainf()
}
