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
		return LoadValueFromText(sc.Text)

	case EvalFile:
		if len(sc.Args) != 1 {
			return cue.Value{}, fmt.Errorf("bad number of args to SourceConfig.File, should be only one filename, got %v", sc.Args)
		}
		return LoadValueFromFile(sc.Args[0])

	case EvalHttp:
		if len(sc.Args) != 1 {
			return cue.Value{}, fmt.Errorf("bad number of args to SourceConfig.Http, should be only one filename, got %v", sc.Args)
		}
		return LoadValueFromHttp(sc.Args[0])

	case EvalBash:
		return LoadValueFromBash(sc.Args)	
	}

	return cue.Value{}, fmt.Errorf("unhandled SourceConfig.Source: %q", sc.Source)
}
