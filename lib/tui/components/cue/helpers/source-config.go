package helpers

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/watch"
	"github.com/hofstadter-io/hof/lib/yagu"
)

type EvalSource string

const (
	EvalNone = ""
	EvalRuntime = "runtime"
	EvalText = "text"
	EvalFile = "file"
	EvalBash = "bash"
	EvalHttp = "http"
	EvalConn = "conn"
)

type SourceConfig struct {
	// manual source data
	Value cue.Value
	Text  string

	// or how to get the value
	Source EvalSource
	Args []string

	_runtime *runtime.Runtime

	ConnGetter func() cue.Value
	// source format here?

	// for handling external changes and updating the Value
	WatchGlobs []string
	WatchTime  time.Duration
	WatchQuit  chan bool
	WatchFunc  func()
}

func (sc *SourceConfig) Encode() (map[string]any, error) {
	return map[string]any{
		"source": sc.Source,
		"args": sc.Args,
		"watch": sc.WatchTime.String(),
	}, nil
}

func (sc *SourceConfig) Decode(input map[string]any) (*SourceConfig, error) {
	aargs := input["args"].([]any)
	args := make([]string, len(aargs))
	for i, a := range aargs {
		args[i] = a.(string)
	}
	watch := "0s"
	if _w, ok := input["watch"]; ok {
		watch = _w.(string)
	}
	d, err := time.ParseDuration(watch)
	if err != nil {
		return nil, fmt.Errorf("error while decoding SourceConfig.watch: %w", err)
	}

	return &SourceConfig{
		Source: EvalSource(input["source"].(string)),
		Args: args,
		WatchTime: d,
	}, nil
}

func (sc *SourceConfig) GetValue() (cue.Value, error) {
	// tui.Log("debug", fmt.Sprintf("SCFG.GetValue %# v", sc))

	switch sc.Source {
	case EvalNone:
		return sc.Value, nil

	case EvalRuntime:
		r, err := LoadRuntime(sc.Args)
		if err != nil {
			return cue.Value{}, err
		}
		sc._runtime = r
		return r.Value, nil

	case EvalText:
		_, v, e := LoadFromText(sc.Text)
		return v, e

	case EvalFile:
		if len(sc.Args) != 1 {
			return cue.Value{}, fmt.Errorf("bad number of args to SourceConfig.File, should be only one filename, got %v", sc.Args)
		}
		_, v, e := LoadFromFile(sc.Args[0])
		return v, e

	case EvalHttp:
		if len(sc.Args) != 1 {
			return cue.Value{}, fmt.Errorf("bad number of args to SourceConfig.Http, should be only one filename, got %v", sc.Args)
		}
		_, v, e := LoadFromHttp(sc.Args[0])
		return v, e

	case EvalBash:
		_, v, e := LoadFromBash(sc.Args)	
		return v, e

	case EvalConn:
		v := sc.ConnGetter()
		return v, nil
	}

	return cue.Value{}, fmt.Errorf("unhandled SourceConfig.Source: %q", sc.Source)
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

func (sc *SourceConfig) Watch(label string, callback func(), debounce time.Duration) error {
	var (
		files []string
		err error
	)
	if len(sc.WatchGlobs) == 0 {
		switch sc.Source {
		case EvalRuntime:
			if sc._runtime == nil {
				r, err := LoadRuntime(sc.Args)
				if err != nil {
					tui.Log("error", err)
					return err
				}
				sc._runtime = r
			}
			files = sc._runtime.GetLoadedFiles()
		case EvalFile:
			files = sc.Args
		default:
			return fmt.Errorf("auto-file discover not available for %s, you can set globs manually though")
		}
	} else {
		files, err = yagu.FilesFromGlobs(sc.WatchGlobs)
	}
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("did not find any files to watch")
	}

	// always kill old watcher
	sc.StopWatch()

	// make a new runner
	sc.WatchQuit = make(chan bool, 2) // non blocking

	cb := func() error {
		callback()
		return nil
	}

	sc.WatchFunc = callback
	err = watch.Watch(cb, files, label, debounce, sc.WatchQuit, false)

	return err
}

func (sc *SourceConfig) StopWatch() {
	if sc.WatchQuit != nil {
		sc.WatchQuit <- true
		sc.WatchQuit = nil
	}
}
