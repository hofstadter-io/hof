package helpers

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/singletons"
)

type EvalSource string

const (
	EvalNone EvalSource = ""
	EvalRuntime EvalSource = "runtime"
	EvalText EvalSource = "text"
	EvalFile EvalSource = "file"
	EvalBash EvalSource = "bash"
	EvalHttp EvalSource = "http"
	EvalConn EvalSource = "conn"
)

type SourceConfig struct {
	Name string

	// manual source data
	Value cue.Value
	Error error
	Text  string

	// or how to get the value
	Source EvalSource
	Args []string

	// Path to fill the CUE value
	// (when used as an input in browser)
	// this will be filled at Path in another value there
	Path string

	_runtime *runtime.Runtime

	ConnGetter func() cue.Value
	// source format here?

	// for handling external changes and updating the Value
	WatchGlobs []string
	WatchTime  time.Duration
	WatchQuit  chan bool
	WatchFunc  func()
}

func (sc *SourceConfig) GetValue() (cue.Value, error) {
	// tui.Log("debug", fmt.Sprintf("SCFG.GetValue %# v", sc))

	switch sc.Source {
	case EvalNone:
		return sc.Value, nil

	case EvalRuntime:
		r, err := LoadRuntime(sc.Args)
		if err != nil {
			return singletons.EmptyValue(), err
		}
		sc._runtime = r
		return r.Value, nil

	case EvalText:
		_, v, e := LoadFromText(sc.Text)
		return v, e

	case EvalFile:
		if len(sc.Args) != 1 {
			return singletons.EmptyValue(), fmt.Errorf("bad number of args to SourceConfig.File, should be only one filename, got %v", sc.Args)
		}
		_, v, e := LoadFromFile(sc.Args[0])
		return v, e

	case EvalHttp:
		if len(sc.Args) != 1 {
			return singletons.EmptyValue(), fmt.Errorf("bad number of args to SourceConfig.Http, should be only one filename, got %v", sc.Args)
		}
		_, v, e := LoadFromHttp(sc.Args[0])
		return v, e

	case EvalBash:
		_, v, e := LoadFromBash(sc.Args)	
		return v, e

	case EvalConn:
		if sc.ConnGetter != nil {
			v := sc.ConnGetter()
			return v, nil
		}
		// we get here on load, which has to decode, and then run connection refresh afterwards
		return singletons.EmptyValue(), nil
	}

	return singletons.EmptyValue(), fmt.Errorf("unhandled SourceConfig.Source: %q", sc.Source)
}

func (sc *SourceConfig) GetText() (string, error) {
	// tui.Log("debug", fmt.Sprintf("SCFG.GetText %# v", sc))
	switch sc.Source {
	case EvalNone:
		return sc.Text, nil

	case EvalRuntime:
		return "", fmt.Errorf("EvalRuntime does not support GetText()")

	case EvalText:
		return sc.Text, nil

	case EvalFile:
		if len(sc.Args) != 1 {
			return "", fmt.Errorf("bad number of args to SourceConfig.File, should be only one filename, got %v", sc.Args)
		}
		s, _, e := LoadFromFile(sc.Args[0])
		return s, e

	case EvalHttp:
		if len(sc.Args) != 1 {
			return "", fmt.Errorf("bad number of args to SourceConfig.Http, should be only one filename, got %v", sc.Args)
		}
		s, _, e := LoadFromHttp(sc.Args[0])
		return s, e

	case EvalBash:
		s, _, err := LoadFromBash(sc.Args)	
		return s, err
	}

	return "", fmt.Errorf("unhandled SourceConfig.Source: %q", sc.Source)
}
