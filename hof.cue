package hof

import (
	g_cuefig "github.com/hofstadter-io/hofmod-cuefig/gen"
	g_cli "github.com/hofstadter-io/hofmod-cli/gen"

	d_cli "github.com/hofstadter-io/hof/design/cli"
	d_cfg "github.com/hofstadter-io/hof/design/config"
)

Cli: _ @gen(cli,hof)
Cli: g_cli.#HofGenerator & {
  Outdir: "./"
  Cli: d_cli.#CLI
}

Context: _ @gen(cuefig,context,ctx)
Context: g_cuefig.#HofGenerator & {
  Outdir: "gen/"
  Config: d_cfg.#HofContext
}

Config: _ @gen(cuefig,config,cfg)
Config: g_cuefig.#HofGenerator & {
  Outdir: "gen/"
  Config: d_cfg.#HofConfig
}

Secret: _ @gen(cuefig,secret,shh,ssh)
Secret: g_cuefig.#HofGenerator & {
  Outdir: "gen/"
  Config: d_cfg.#HofSecret
}

UserContext: _ @gen(cuefig,ucontext,uctx)
UserContext: g_cuefig.#HofGenerator & {
  Outdir: "gen/"
  Config: d_cfg.#HofUserContext
}

UserConfig: _ @gen(cuefig,uconfig,ucfg)
UserConfig: g_cuefig.#HofGenerator & {
  Outdir: "gen/"
  Config: d_cfg.#HofUserConfig
}

UserSecret: _ @gen(cuefig,usecret,ushh,ussh)
UserSecret: g_cuefig.#HofGenerator & {
  Outdir: "gen/"
  Config: d_cfg.#HofUserSecret
}

Foo: _ @gen(foo,bar=fuck)

