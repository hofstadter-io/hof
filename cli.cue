package hof

import (
	"github.com/hofstadter-io/hofmod-cli:cli"
	"github.com/hofstadter-io/hof/design"
)

Outdir :: "./"

HofGenCli: cli.HofGenerator & {
  Outdir: "./"
  Cli: design.CLI
}

