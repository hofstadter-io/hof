package lang

import (
	"bytes"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/hofstadter-io/hof/pkg/lang/context"
	"github.com/hofstadter-io/hof/pkg/lang/phases"
	"github.com/hofstadter-io/hof/pkg/lang/visit"
)

var Phases = []phases.Phase {
	phases.LoadModule,
	phases.CheckPackages,
	phases.FindDefinitions,
	phases.FillScopes,
	phases.ScopeDebug,
}

func Eval(entrypoint string) error {
	ctx := context.NewContext()
	ctx.Entrypoint = entrypoint

	for _, phase := range Phases {
		phase(ctx)

		if len(ctx.Errors) > 0 {
			ctx.PrintErrors()
			return errors.New("Failed to load")
		}
	}

	data, err := visit.ToData(ctx.Module)
	if err != nil {
		return err
	}

	return nil

	err = dump(data)
	if err != nil {
		return err
	}

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
	fmt.Println(b.String())
	return nil
}

