package helpers

import (
	"fmt"
	"strings"
	"time"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/tui"
)

func CreateFrom(args []string, context map[string]any) (*SourceConfig, error) {
	tui.Log("extra", fmt.Sprintf("SourceConfig.CreateFrom: %v %#v", args, context))

	sc := new(SourceConfig)
	if len(args) > 0 && strings.HasPrefix(args[0], "@") {
		sc.Path = args[0][1:]
		args = args[1:]
	}

	// source, defaults to CUE runtime / loader
	if _src, ok := context["source"]; ok {
		sc.Source = EvalSource(_src.(string))
	} else {
		sc.Source = EvalRuntime
	}

	// handle extra settings per EvalSource type
	switch sc.Source {
	case EvalConn:
	  if _fn, ok := context["fn"]; ok {
			fn := _fn.(func () cue.Value)
			sc.ConnGetter = func() cue.Value {
				return fn()
			}
		}
	}

	sc.Args = args 
	return sc, nil
}

func (sc *SourceConfig) UpdateFrom(action string, args []string, context map[string]any) (bool, error) {
	var err error
	handled := true

	switch action {
	case "set":
		if len(args) > 0 && strings.HasPrefix(args[0], "@") {
			sc.Path = args[0][1:]
			args = args[1:]
		}
		sc.Source = EvalSource(context["source"].(string))
		sc.Args = args

		// TODO, handle connections, how? maybe set from outside manually for now?
		// or maybe we should turn source config into an interface ( and package definitely )

	case "globs":
		sc.WatchGlobs = args

	case "watch":
		// parse duration if available
		d := time.Duration(42*time.Millisecond)
		if len(args) < 1 {
			//aerr := fmt.Errorf("watch requires a duration like 1s or 300ms")
			//tui.Log("warn", aerr)
			tui.Log("warn", fmt.Sprintf("no watch duration given, setting to %s", d))
			tui.StatusMessage(fmt.Sprintf("no watch duration given, setting to %s", d))
		} else {
			d, err = time.ParseDuration(args[0])
		}
		sc.WatchTime = d

	default:
		handled = false

	}


	return handled, err
}

func extractInfo(args []string, context map[string]any) {
	//source := ""
	//if _source, ok := context["source"]; ok {
	//  source = _source.(string)
	//}

	//target := "value"
	//if _target, ok := context["target"]; ok {
	//  target = _target.(string)
	//}

}

