package hof

import (
	cuefig "github.com/hofstadter-io/hofmod-cuefig/gen"
	cli "github.com/hofstadter-io/hofmod-cli/gen"
	"github.com/hofstadter-io/hof/design"
)

HofGenCli: cli.#HofGenerator & {
  Outdir: "./"
  Cli: design.#CLI
}

HofGenConfig: cuefig.#HofGenerator & {
  Outdir: "./"
  Config: design.#HofConfig
}

HofGenCreds: cuefig.#HofGenerator & {
  Outdir: "./"
  Config: design.#HofCredentials
}
