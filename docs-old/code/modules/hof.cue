package hof

import (
	// CUE stdlib
	"strings"

	// dependency import
	"github.com/hofstadter-io/hofmod-cli/gen"

	// module-local import
	"github.com/hofstadter-io/hof/design"
)

Cli: gen.#Generator & {
	@gen(cli,hof)
	Outdir: "./"
	Cli:    design.#CLI
	WatchGlobs: ["./design/**/*"]
}
