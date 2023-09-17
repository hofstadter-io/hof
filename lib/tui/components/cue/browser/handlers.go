package browser

import (
	"fmt"
	"strconv"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
)

func handleSet(B *Browser, action string, args []string, context map[string]any) (bool, error) {
	sc, err := helpers.CreateFrom(args, context)
	if err != nil {
		return true, err
	}
	B.sources = []*helpers.SourceConfig{sc}

	B.setThinking(true)
	defer B.setThinking(false)
	B.RebuildValue()
	B.Rebuild()

	return true, nil
}

func handleAdd(B *Browser, action string, args []string, context map[string]any) (bool, error) {
	sc, err := helpers.CreateFrom(args, context)
	if err != nil {
		tui.Log("alert", fmt.Sprint("GOT HERE: ", err))
		return true, err
	}
	B.sources = append(B.sources, sc)

	B.setThinking(true)
	defer B.setThinking(false)
	B.RebuildValue()
	B.Rebuild()

	return true, nil
}

func handleWatchConfig(B *Browser, action string, args []string, context map[string]any) (bool, error) {

	do := func(cfg *helpers.SourceConfig, action string, args []string, context map[string]any) (bool, error) {
		h, err := cfg.UpdateFrom(action, args, context)
		if err != nil {
			return h, err
		}
		cfg.WatchFunc = func() {
			B.setThinking(true)
			defer B.setThinking(false)
			B.RebuildValue()
			B.Rebuild()
		}
		go cfg.Watch()
		return true, nil
	}

	return B.handleOneOrAll(do, action, args, context)
}

type doer func(cfg *helpers.SourceConfig, action string, args []string, context map[string]any) (bool, error)

func (B *Browser) handleOneOrAll(do doer, action string, args []string, context map[string]any) (bool, error) {
	var err error
	var handled bool

	// look for an index marker (2 args)
	idx := -1
	if len(args) > 1 {
		// there was an index setting
		idxStr := args[0]
		args = args[1:]
		idx, err = strconv.Atoi(idxStr)
	}
	if err != nil {
		return true, err
	}

	if idx > -1 {
		cfg := B.sources[idx]
		handled, err = do(cfg, action, args, context)
	} else {
		for _, cfg := range B.sources {
			handled, err = do(cfg, action, args, context)
			if err != nil {
				break
			}
		}
	}

	return handled, err
}
