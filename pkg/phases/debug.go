package phases

import (
  "fmt"

	"github.com/hofstadter-io/hof/pkg/context"
)

func Debug(ctx *context.Context) error {
	for _, pkg := range ctx.Packages {
    fmt.Println(pkg.Path)
	}

	return nil
}


