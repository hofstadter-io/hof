package hof

import (
	"github.com/hofstadter-io/hofmod-cli/gen"

	"github.com/hofstadter-io/hof/design"
)

Cli: gen.Generator & {
	@gen(cli,hof)
	Outdir: "./"
	Cli:    design.CLI
	WatchGlobs: ["./design/**/*"]
	WatchExec: {
		@task(os.Exec)
		cmd: ["go", "install", "\(Outdir)/cmd/hof"]
		env: {
			CGO_ENABLE: "0"
		}
		exitcode: _
	}
}
