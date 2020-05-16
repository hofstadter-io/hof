package hof

import (
	g_cuefig "github.com/hofstadter-io/hofmod-cuefig/gen"
	g_cli "github.com/hofstadter-io/hofmod-cli/gen"

	d_cli "github.com/hofstadter-io/hof/design/cli"
	d_cfg "github.com/hofstadter-io/hof/design/config"
)

HofGenCli: g_cli.#HofGenerator & {
  Outdir: "./"
  Cli: d_cli.#CLI
}

HofGenConfig: g_cuefig.#HofGenerator & {
  Outdir: "./"
  Config: d_cfg.#HofConfig
}

HofGenCreds: g_cuefig.#HofGenerator & {
  Outdir: "./"
  Config: d_cfg.#HofCredentials
}
