package extension

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers"
)

func Run(args []string, rflags flags.RootPflagpole) error {
	r, err := runtime.NewRuntime()
	if err != nil {
		return fmt.Errorf("failed to create runtime: %v", err)
	}

	handlers.SetupHandlers(r)

	return r.Run()
}
