package hof

import (
	cli "github.com/hofstadter-io/hofmod-cli/gen"
	"github.com/hofstadter-io/hof/design"
)

HofGenCli: cli.#HofGenerator & {
  Outdir: "./"
  Cli: design.#CLI
}

