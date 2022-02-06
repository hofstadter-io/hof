// This package provides a context for Tasks
// and a registry for their usage in flows.
package context

import (
	"context"
	"io"
	"sync"

	"cuelang.org/go/cue"
)

// A Context provides context for running a task.
type Context struct {
	Context context.Context
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Value   cue.Value
	Error   error

  // Global lock around CUE evaluator 
  CUELock  *sync.Mutex

  // map of cue.Values
  ValStore *sync.Map

  // map of chan?
  Mailbox  *sync.Map

  // channels for
  // - stats & progress

  // debug / internal
  DebugTasks bool
  Verbosity  int
}

// consider adding here... a
// global registry of named channels

// A RunnerFunc creates a Runner.
type RunnerFunc func(v cue.Value) (Runner, error)

// A Runner defines a task type.
type Runner interface {
	// Runner runs given the current value and returns a new value which is to
	// be unified with the original result.
	Run(ctx *Context) (results interface{}, err error)
}

// Register registers a task for cue commands.
func Register(key string, f RunnerFunc) {
	runners.Store(key, f)
}

// Lookup returns the RunnerFunc for a key.
func Lookup(key string) RunnerFunc {
	v, ok := runners.Load(key)
	if !ok {
		return nil
	}
	return v.(RunnerFunc)
}

var runners sync.Map
