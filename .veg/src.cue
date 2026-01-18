package veg

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

src: {
	repo: env.#GitRepo & {
		@env(src-hof-repo)

		// name: string | *"repo"
		url: flags.repo
		ref: string | *flags.ref
	}
	disk: env.#HostDir & {
		@env(src-hof-disk)

		// name: string | *"disk"
		path:      flags.disk
		gitignore: true
	}
	adk: {
		disk: env.#HostDir & {@env(src-adk-disk), path: flags.adk}
		fork: env.#GitRepo & {@env(src-adk-fork), url: "https://github.com/verdverm/adk-go", ref: "veg"}
	}
	dagger: {
		disk: env.#HostDir & {@env(src-dagger-disk), path: flags.dagger}
		fork: env.#GitRepo & {@env(src-dagger-fork), url: "https://github.com/verdverm/dagger", ref: "patches"}
	}

	// setup code base on flags and value
	code: {@env(src-hof-curr), name: "code"}
	if flags.src == "repo" {code: repo}
	if flags.src == "disk" {code: disk}
	if strings.HasPrefix(flags.src, "https://") {
		code: env.#GitRepo & {url: flags.src}
	}
	if strings.HasPrefix(flags.src, ".") {
		code: env.#HostDir & {path: flags.src}
	}

	cuemod: env.#Dir & {
		@env(src-cuemod)
		sources: [src.code]
		include: [
			"cue.mod/module.cue",
			// "*.cue", // eventually, when we rework all of ci, use .veg more, and have a root index that imports many things, like a mega package if the user wants
			"schemas",
			"catalogs/env",
			"examples/env",
			"flow/tasks/*.cue",
			"flow/tasks/*/*.cue",
			"lib/env/common",
			"SECURITY.md",
			"README.md",
			"AGENTS.md",
			"LICENSE",
		]
	}

	extn: {
		vscode: {}
	}
}

host: {
	docker: {
		// get a socket from the host, inception will be painfully slow otherwise
		socket: env.#HostSocket & {
			@env(docker-sock)
			name: "docker-sock"
			path: flags.socket
		}
	}

	kindapi: {
		service: env.#HostService & {
			@env(kindapi-host)
			name: "kindapi-host"
			host: "host.docker.internal"
			ports: [{port: 6443}]
		}
	}
}

secrets: {
	dotssh: env.#HostDir & {
		@env(shh-dotssh)
		path: "~/.ssh"
	}
	kubecfg: env.#HostDir & {
		@env(shh-kubecfg)
		path: "~/.kube"
	}

	google: env.#Secret & {
		@env(shh-google)
		name:   "google-api-key"
		source: "GOOGLE_API_KEY"
	}
}
