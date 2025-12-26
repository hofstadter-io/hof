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
			_steps.apt & {#pkgs: [
				"ca-certificates",
				"curl",
				"gnupg",
				"git",
				"make",
				"unzip",
				"wget",
				"xz-utils",
				"zsh",
			]},

			// term customization
			_tools.zsh,
			env.Args & {args: ["zsh"]},
			env.Term & {args: ["zsh"]},
			env.Entrypoint & {args: ["zsh"]},

      // default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},
		]
	}

	dev: env.Container & {
		@env()
		name: "veg-dev"
		#hof: description: "A development image with many tools"
		from: base

		steps: [
      // todo, put these with the tools that depend on them (if they are one)
			_steps.apt & {#pkgs: [
        // go
				"g++",
				"gcc",
				"libc6-dev",
				"netbase",
				"pkg-config",
				"sq",

        // python
        "pip",
        "pipx",
        "pylint",
        "python3-poetry",
        "python3-pytest",
        "python3-flake8",
			]},

			// setup languages
			_tools.go,
			_tools.node,
			_tools.python,

			// other binary tools
			_tools.githubBin & {#repo: "cue-lang/cue", #ver: "0.15.1"},

      // tools for agents
      _tools.agents.lsp2mcp,
		]
	}

	ops: env.Container & {
		@env()
		name: "veg-ops"
		#hof: description: "Extension to veg-dev to add devops tooling"
		from: dev

		steps: [
			_tools.kubectl,
			_tools.helm,
			_tools.hashicorpBin & {#tool: "terraform"},
			_tools.hashicorpBin & {#tool: "packer"},
      _tools.crane,
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

	vegeta: env.Container & {
		@env()
		name: "vegeta"
		#hof: description: "all of the veggie images, it's over 9000"
		from: ops

    // TODO, realize the dagger way of diamond build pattern optimizations (once we branch more than 2 wide, and have _tools in a better place with #file/#dir)
    steps: incept.steps

  }

}

// dind images (docker & dagger)
// bin tools for File / Dir / Copy
// - python 3, poetry, uv
// - gcloud
//
// build hof & formatters
// docker-compose like experience (testnet?)
//