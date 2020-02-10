package phases

import (
	"errors"
	"os"

	"github.com/hofstadter-io/hof/pkg/context"
)

func LoadModule(ctx *context.Context) error {
	info, err := os.Lstat(ctx.Entrypoint)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New("Entrypoint should be a directory")
	}

	_, err = ctx.LoadModule(ctx.Entrypoint)
	if err != nil {
		return err
	}

	return nil
}
