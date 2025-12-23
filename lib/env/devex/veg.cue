package devex

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

base: env.Container & {
	@env()
	#hof: description: " base image for veg"
	from: "debian:13-slim"

	_pkgs: [
		"curl",
		"git",
		"zsh",
	]

	steps: [
    
    // basics
		env.Workdir & { path: "/root" },
		_steps.apt & {#pkgs: _pkgs},

	]

}

veg: env.Container & {
	@env()
	#hof: description: " base image for veg"
	from: base

  steps:[
    _tools.zsh,
		env.Term & {args: ["zsh"]},
  ]

}
// registry:2 as service & publishing
// dind images (docker & dagger)
// bin tools for File / Dir / Copy
// - zsh, omzsh
// - go, gopls
// - python 3, poetry, uv
// - node, pnpm
// - cue, dagger, docker
// - helm, tf, k8s
// - gcloud
//
// ////
//
// build hof & formatters
// docker-compose like experience (testnet?)
//
//
