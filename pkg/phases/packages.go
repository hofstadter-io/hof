package phases

import (
	"github.com/hofstadter-io/hof/pkg/context"
	"github.com/hofstadter-io/hof/pkg/visit"
)

func CheckPackages(ctx *context.Context) error {
	for _, pkg := range ctx.Packages {
		err := visit.CheckPackageFileDefines(pkg)
		if err != nil {
			ctx.AddError(err)
		}
	}

	return nil
}
