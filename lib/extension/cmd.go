package extension

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/agent"
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

	// fmt.Println("R.Agentics:", len(r.Agentics))

	err = r.InitServices()
	if err != nil {
		return fmt.Errorf("failed to init services: %v", err)
	}

	err = r.EnrichAgentic(nil, AgenticEnricher)
	if err != nil {
		return err
	}

	ar, err := aruntime.NewRuntime(r)
	if err != nil {
		return fmt.Errorf("failed to create agent runtime: %v", err)
	}
	// fmt.Println("AR.Agentics:", len(ar.Agentics))

	ar.BackfillAgentic()

	// fmt.Println("BR.Agentics:", len(ar.Agentics))
	ws.SetupHandlers(r, ar)

	return ar.Run()
}

func AgenticEnricher(R *runtime.Runtime, e *agent.Agentic) error {
	// no-op for now
	return nil
}
