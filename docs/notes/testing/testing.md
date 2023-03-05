---
title: "Testing"
weight: 18
---

__hof__ has several testings systems powered by CUE.
This enables far more reuse of
snippets, test cases, and commands.
Tests can also be selected by attribute
for targeted execution.
You also get all the power and flexibility of CUE
when defining tests and cases.

[Running tests](#running-tests-with-hof) with `hof test`
enables you to setup and run any tests
using Cue attributes and CLI flags.

The types of tests are:

- [External](#external-testing) test systems can be invoked.
- [TSuite](#tsuite-testing) for Cue based test cases.
- [API](#api-testing) for testing REST and GraphQL servers.
- [HLS](#hls-testing) for more complex, text file based setups.

This section has an introduction to __hof's__ testing facilities.
You can learn more in the [testing](/testing/) section.



## Running tests with __hof__

This is __hof__'s current test file.
You can first see the setup for Golang based testers.
Following that are a series of `@test(...)` suites and tests.
`@test(suite,...)` is a suite. Testers have their base type
(bash,api,...) as the first arg in the attribute.

#### `test.cue`:

```text
package hof

import "strings"

// Defined partial test configurations

#GoBaseTest: {
	skip: bool | *false

	sysenv: bool | *false
	env?: [string]: string

	dir: string
	...
}

// Specialize base for go tests
#GoBashTest: #GoBaseTest & {
  // ... omit for brevity, see below
}

// Specialize base for go coverage report generation
#GoBashCover: #GoBaseTest & {
  // ... omit for brevity, see below
}

// Actual test configuration

// Test hof generated code
gen: _ @test(suite,gen)
gen: {
	cmds: #GoBashTest @test(bash,test,cmd)
  // ... omit for brevity, see below
}

// Test Hof Linear Script (hls)
hls: _ @test(suite,hls)
hls: {
	runtime: #GoBashTest @test(bash,test,runtime)
	runtime: #GoBashTest @test(bash,test,runtime)
  // ... omit for brevity, see below
}

// Test Hof libraries
lib: _ @test(suite,lib)
lib: {
	mod: #GoBashTest @test(bash,test,mod)
	modC: #GoBashCover @test(bash,cover,mod)
  // ... omit for brevity, see below
}
```

Tests are run with the __hof__ CLI.

```sh
hof test [-s <suite>] [-t <test>] <test.cue>
```

Each flag can be specified multiple times.
The flag values support regex, so "st" will match 'st' and 'test', so use '^st$' for just st

- Multiple `-s` suite flags "or" together
- Multiple `-t` tester flags "and" together


```sh
# Run all tests
hof test

# Run all lib tests
hof test -s lib

# Run all cover tests
hof test -t cover

# Run lib/st tests
hof test -s lib -t test -t "^st$"
```

The code implementing test functionality can be found here:
https://github.com/hofstadter-io/hof/tree/_dev/lib/test


## External testing

Most languages come with robust testing systems or libraries.
__hof__ does not try to replace these, rather make them easier
to invoke, especially in a polyglot or monorepo setup.

```text
#GoBaseTest: {
	skip: bool | *false

	sysenv: bool | *false
	env?: [string]: string

	dir: string
	...
}

#GoBashTest: #GoBaseTest & {
	dir: string
	script: string | *"""
	rm -rf .workdir
	go test -cover ./
	"""
	...
}

#GoBashCover: #GoBaseTest & {
	dir: string
	back: strings.Repeat("../", strings.Count(dir, "/") + 1)
	script: string | *"""
	rm -rf .workdir
	go test -cover ./ -coverprofile cover.out -json > tests.json
	"""
	...
}
```

This is the eaxmple usage:

```text
lib: _ @test(suite,lib)
lib: {

	mod: #GoBashTest @test(bash,test,mod)
	mod: {
		dir: "lib/mod"
	}
	modC: #GoBashCover @test(bash,cover,mod)
	modC: {
		dir: "lib/mod"
	}

	st: #GoBashTest @test(bash,test,st)
	st: {
		dir: "lib/structural"
	}
	stC: #GoBashCover @test(bash,cover,st)
	stC: {
		dir: "lib/structural"
	}

}
```



## TSuite testing

Cue's ability to express data and configuration
can also be used to express test cases.
You will need to setup:

- A Go test file for the drive
- Cue files for the test cases

While the driver file is written in Go, you are not limited to testing just Go.
For Go, you setup the function you wish to test and the directory containing test cases.
When testing other languages, the process involves
a step to convert the Cue test cases to the input for the target
and using the previous external invocation from the previous section.


#### driver_test.go:

```go
package structural_test

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

var (
  // for printing
	MyFuncFmtStr = "MyFunc[%v]: %v"
  // key to lookup cases in Cue code
	MyFuncTestCases = []string{
		"#MyFuncCases",
	}
)

// Define a TSuite
type MyTestSuite struct {
	*cuetils.TestSuite
}

// New func for TSuite
func NewMyTestSuite() *MyTestSuite {
	ts := cuetils.NewTestSuite(nil, MyOp)
	return &MyTestSuite{ ts }
}

// How TSuite finds things to run, this function name fits a pattern
func TestMyTestSuites(t *testing.T) {
	suite.Run(t, NewMyTestSuite())
}

// An op to wrap the function we wish to run, looking up the keys from the args
func MyOp(name string, args cue.Value) (val cue.Value, err error) {
	orig := args.Lookup("orig")
	next := args.Lookup("next")
	return structural.DiffValues(orig, next)
}

// Function which setups and runs all the cases
func (PTS *MyTestSuite) TestMyFuncCases() {

	err := PTS.SetupCue()
	assert.Nil(PTS.T(), err, fmt.Sprintf(MyFuncFmtStr, "setup", "Loading test cases should return non-nil error"))
	if err != nil {
		return
	}

	tSyn, err := cuetils.ValueToSyntaxString(PTS.CRT.CueValue)
	assert.Nil(PTS.T(), err, fmt.Sprintf(MyFuncFmtStr, "syntax", "Printing test cases should return non-nil error"))
	if err != nil {
		fmt.Println(tSyn)
		return
	}

	PTS.Op = MyOp
	PTS.RunCases(MyFuncTestCases)
}

```

#### cases.cue:

```text
package testdata

#MyFuncCases: {
  // test group
	simple: @group(simple)
	simple: {
    // test case
		t_0001: {
      // args to function
			args: {
        ...
			}
      // expected result
			ex: {
        ...
			}
		}
		t_0002: {
      // another case...
		}
	}
  // another test group
  complex: @group(complex)
	complex: {
    ...
	}
}
```

You have quite a bit of flexibility for the contents of the CUE files.
It will depend on what you are testing and its arguments and returns.
Attributes are used so you can run specific groups or cases with granularity.


## API testing

Similar to TSuite testing, you can setup your test cases and configuration
using the power of Cue and __hof__.
The difference here is that you only need to write Cue files
and can omit language specific files.

More to come soon!



## HLS testing

__hof__ has a scripting system called
Hofstadter Linear Script (HLS).
HLS is designed to combine

- the direct executability of Bash
- ability to work with data objects from Cue+__hof__
- test case orientedness from Golang's testsuite
- methods for working with multiple files in a single file from Golang's txtar

This setup allows you to write tests as text files, more specifically as HLS scripts.
Rather than writing and compiling code, you can add test cases to a directory
while providing a common setup and environment.

The following is taken from __hof's__ `lib/mod/cli_test.cue`:

```go
package mod_test

import (
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		env.Vars = append(env.Vars, "GITHUB_TOKEN="+token)
	}
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	return nil
}

func TestModTests(t *testing.T) {
	yagu.Mkdir(".workdir/tests")
	runtime.Run(t, runtime.Params{
		Setup: envSetup,
		Dir: "testdata",
		Glob: "*.txt",
		WorkdirRoot: ".workdir/tests",
	})
}

func TestModBugs(t *testing.T) {
	yagu.Mkdir(".workdir/bugs")
	runtime.Run(t, runtime.Params{
		Setup: envSetup,
		Dir: "testdata/bugs",
		Glob: "*.txt",
		WorkdirRoot: ".workdir/bugs",
	})
}
```

HLS scripting is a subject of its own.
You can learn more in the [HLS scripting section](/hls-scripting/).

