package pkg

import (
	"bytes"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/hofstadter-io/hof/pkg/context"
	"github.com/hofstadter-io/hof/pkg/phases"
)

var Phases = []phases.Phase {
	phases.LoadModule,
	phases.CheckPackages,
	phases.FindDefinitions,
	// phases.Debug,
}

func Do(entrypoint string) error {
	ctx := context.NewContext()
	ctx.Entrypoint = entrypoint

	for _, phase := range Phases {
		phase(ctx)

		if len(ctx.Errors) > 0 {
			ctx.PrintErrors()
			return errors.New("Failed to load")
		}
	}

	/*
	err := dump(ctx.Module)
	if err != nil {
		return err
	}
	*/

	return nil
}

func dump(thing interface{}) error {
	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)

	err := encoder.Encode(thing)
	if err != nil {
		return err
	}

	fmt.Println("============================================")
	fmt.Printf("Context:\n%s\n", b.String())
	return nil
}

