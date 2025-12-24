package devex

import (
	"github.com/hofstadter-io/hof/schemas/env"
)


veg: {
  base: env.Container & {
    @env()
    name: "veg-base"
    #hof: description: "A minimal debian image with a few common tools"
    from: "debian:13-slim"

    _pkgs: [
      "curl",
      "git",
    ]

    steps: [
      
      // basics
      env.Workdir & { path: "/root" },
      _steps.apt & {#pkgs: _pkgs},

    ]

  }

  dev: env.Container & {
    @env()
    name: "veg-dev"
    #hof: description: "A development image with many tools"
    from: base

    steps:[
      _steps.apt & {#pkgs: ["zsh"]},
      _tools.zsh,
      env.Term & {args: ["zsh"]},
    ]
  }

  incept: env.Container & {
    @env()
    name: "veg-incept"
    #hof: description: "Extension to veg-dev to add docker/dagger setup for inception"
    from: dev

    steps:[]
  }


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
