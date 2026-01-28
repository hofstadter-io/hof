package extension

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	aruntime "github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/ws"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func Run(args []string, rflags flags.RootPflagpole) error {
	// create our core runtime
	r, err := runtime.New(args, rflags)
	if err != nil {
		return fmt.Errorf("failed to create veg runtime: %v", err)
	}

	err = r.Load()
	if err != nil {
		return cuetils.ExpandCueError(err)
	}

	err = r.InitServices()
	if err != nil {
		return fmt.Errorf("failed to init services: %v", err)
	}

	ar, err := aruntime.NewRuntime(
		r.DB,
		r.Envs,
		r.Agentics,
	)
	if err != nil {
		return fmt.Errorf("failed to create agent runtime: %v", err)
	}
	ar.BackfillAgentic()

	ws.SetupHandlers(ar)

	return ar.Run()
}
