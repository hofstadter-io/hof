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

		steps: [
			// basics
			env.Workdir & {path: "/root"},
			_steps.apt & {#pkgs: [
				"ca-certificates",
				"curl",
				"gnupg",
				"git",
				"make",
				"wget",
			]},
		]
	}

	dev: env.Container & {
		@env()
		name: "veg-dev"
		#hof: description: "A development image with many tools"
		from: base

		steps: [
			_steps.apt & {#pkgs: [
				"g++",
				"gcc",
				"libc6-dev",
				"netbase",
				"sq",
				"pkg-config",
				"unzip",
				"xz-utils",
				"zsh",
			]},

			// setup languages
			_tools.go & {#ver: "1.25.4"},
			env.Env & {PATH: "$PATH:/usr/local/go/bin"},
			_tools.node & {#ver: "24.12.0"},

			// other binary tools
			_tools.githubBin & {#repo: "cue-lang/cue", #ver: "0.15.1"},
			_tools.kubectl & {#ver: "1.31.0"},
			_tools.helm & {#ver: "4.0.4"},
			_tools.hashicorpBin & {#tool: "terraform", #ver: "1.14.3"},
			_tools.hashicorpBin & {#tool: "packer", #ver: "1.14.3"},
      _tools.crane & { #ver: "0.20.7" },

			// term customization
			_tools.zsh,
			env.Term & {args: ["zsh"]},
		]
	}

	incept: env.Container & {
		@env()
		name: "veg-incept"
		#hof: description: "Extension to veg-dev to add docker/dagger setup for inception"
		from: dev

		steps: [
			_tools.githubBin & {#repo: "dagger/dagger", #ver: "0.19.7"},
			// todo, nested steps
			_tools.dockerRepo,
			_steps.apt & {#pkgs: ["docker-ce-cli", "docker-buildx-plugin", "docker-compose-plugin"]},
		]

		// whatever we import / user here, should also have mounts defined for easy reuse for runtime (run/up/asService)
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
