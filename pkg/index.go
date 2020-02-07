package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/hofstadter-io/hof/pkg/context"
	// "github.com/hofstadter-io/hof/pkg/walkers"
)

func Do(entrypoint string) error {
	info, err := os.Lstat(entrypoint)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New("Entrypoint should be a directory")
	}

	ctx := context.NewContext()

	_, err = ctx.LoadModule(entrypoint)
	if err != nil {
		return err
	}

	if len(ctx.Errors) > 0 {
		hadErrors := ctx.PrintErrors()
		if hadErrors {
			return errors.New("Failed to load")
		}
	}
	// ctx.Print()

	for pname, pkg := range ctx.Packages {
		fmt.Println("Walking", pname)
		for fname, _ := range pkg.Files {
			fmt.Println(" -", fname)
			// walkers.Print(file)
		}
	}

	/*
	err = dump(context)
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

	fmt.Printf("Context:\n%s\n", b.String())
	return nil
}

