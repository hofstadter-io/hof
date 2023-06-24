package runtime

import (
	"errors"
	"fmt"
)

var (
	failedRun = errors.New("failed run")
	skipRun   = errors.New("skip")
)

type Runner struct {
	LogLevel string

	Failed bool
}

func (r *Runner) Skip(is ...interface{}) {
	// panic(skipRun)
}

func (r *Runner) Fatal(is ...interface{}) {
	r.Log(is...)
	r.Failed = true
	// r.FailNow()
}

func (r *Runner) Parallel() {
	// No-op for now; we are currently only running a single script in a
	// testscript instance.
}

func (r *Runner) Log(is ...interface{}) {
	fmt.Print(is...)
}

func (r *Runner) FailNow() {
	panic(failedRun)
}

func (r *Runner) Run(n string, f func(T)) {
	// For now we we don't top/tail the run of a subtest. We are currently only
	// running a single script in a cript instance, which means that we
	// will only have a single subtest.
	f(r)
}

func (r *Runner) Verbose() bool {
	return r.LogLevel != ""
}
