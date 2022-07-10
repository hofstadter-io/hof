package hof

import (
	g_cuefig "github.com/hofstadter-io/hofmod-cuefig/gen"
	g_cli "github.com/hofstadter-io/hofmod-cli/gen"

	d_cli "github.com/hofstadter-io/hof/design/cli"
	d_cfg "github.com/hofstadter-io/hof/design/config"
)

Cli: g_cli.#HofGenerator & {
	@gen(cli,hof)
	Outdir: "./"
	Cli:    d_cli.#CLI
	WatchGlobs: ["./design/**/*"]
	WatchXcue: ["./cue.mod/**/*"]
}

Context: g_cuefig.#HofGenerator & {
	@gen(cuefig,context,ctx)
	Outdir: "gen/"
	Config: d_cfg.#HofContext
}

Config: g_cuefig.#HofGenerator & {
	@gen(cuefig,config,cfg)
	Outdir: "gen/"
	Config: d_cfg.#HofConfig
}

Secret: g_cuefig.#HofGenerator & {
	@gen(cuefig,secret,shh,ssh)
	Outdir: "gen/"
	Config: d_cfg.#HofSecret
}

UserContext: g_cuefig.#HofGenerator & {
	@gen(cuefig,ucontext,uctx)
	Outdir: "gen/"
	Config: d_cfg.#HofUserContext
}

UserConfig: g_cuefig.#HofGenerator & {
	@gen(cuefig,uconfig,ucfg)
	Outdir: "gen/"
	Config: d_cfg.#HofUserConfig
}

UserSecret: g_cuefig.#HofGenerator & {
	@gen(cuefig,usecret,ushh,ussh)
	Outdir: "gen/"
	Config: d_cfg.#HofUserSecret
}
