package hof

import (
	"github.com/hofstadter-io/hofmod-cli/gen"

	"github.com/hofstadter-io/hof/design"
)

Cli: gen.#Generator & {
	@gen(cli,hof)
	Outdir: "./"
	Cli:    design.#CLI
	WatchGlobs: ["./design/**/*"]
}
