package script

import (
	"fmt"
)

type Tester struct {
	verbose bool
}

func (t Tester) Skip(is ...interface{}) {
	panic(skipRun)
}

func (t Tester) Fatal(is ...interface{}) {
	t.Log(is...)
	t.FailNow()
}

func (t Tester) Parallel() {
	// No-op for now; we are currently only running a single script in a
	// testscript instance.
}

func (t Tester) Log(is ...interface{}) {
	fmt.Print(is...)
}

func (t Tester) FailNow() {
	panic(failedRun)
}

func (t Tester) Run(n string, f func(T)) {
	// For now we we don't top/tail the run of a subtest. We are currently only
	// running a single script in a cript instance, which means that we
	// will only have a single subtest.
	f(t)
}

func (t Tester) Verbose() bool {
	return t.verbose
}

