package helpers

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"
)

type EvalSource string

const (
	EvalNone = ""
	EvalRuntime = "runtime"
	EvalText = "text"
	EvalFile = "file"
	EvalBash = "bash"
	EvalHttp = "http"
)

type SourceConfig struct {
	// manual source data
	Value cue.Value
	Text  string

	// or how to get the value
	Source EvalSource
	Args []string
	Watch bool
	Refresh time.Duration

	// source format here?
}

func (sc SourceConfig) Encode() (map[string]any, error) {
	return map[string]any{
		"source": sc.Source,
		"args": sc.Args,
		"watch": sc.Watch,
		"refresh": sc.Refresh.String(),
	}, nil
}


func (sc SourceConfig) GetValue() (cue.Value, error) {
	// tui.Log("debug", fmt.Sprintf("SCFG.GetValue %# v", sc))

	switch sc.Source {
	case EvalNone:
		return sc.Value, nil

	case EvalRuntime:
		r, err := LoadRuntime(sc.Args)
		if err != nil {
			return cue.Value{}, err
		}
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
	}

	return cue.Value{}, fmt.Errorf("unhandled SourceConfig.Source: %q", sc.Source)
}

func (sc SourceConfig) GetText() (string, error) {
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
