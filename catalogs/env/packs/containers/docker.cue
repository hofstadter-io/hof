package containers

import "github.com/hofstadter-io/hof/schemas/env"

docker: {
	#ver: string | *"29"

	img: {
		base: env.#Container & {
			@env(docker-img-base)
			from: "docker:\(#ver)"
		}
		cli: env.#Container & {
			@env(docker-img-cli)
			from: "docker:\(#ver)-cli"
		}
		dind: env.#Container & {
			@env(docker-img-dind)
			from: "docker:\(#ver)-dind"
		}
	}

	cli: {
		binary: env.#File & {
			@env(docker-cli)
			path:   "/usr/local/bin/docker"
			source: img.cli
		}
		install: env.File & {
			path:    "/usr/local/bin/docker"
			content: cli.binary
		}
	}

	daemon: {
		ctr: env.#Container & {
			from: img.dind
			steps: [
				env.Mount & {path: "/tmp", source: vols.tmp},
				// this seems bad, we probably don't really want to use this, but just the client and host socket mount
				// env.Mount & {path: "/var/lib/docker", source: vols.varlib},
			]
		}
		svc: env.#Service & {
			@env(docker-daemon)
			hostname: "docker-daemon"
			source:   daemon.ctr
			ports: [{port: 2375}]
			insecureRootCapabilities: true
		}

		bind: [
			env.EnvVars & {DOCKER_HOST: "tcp://docker-daemon:2375"},
			env.BindService & {service: daemon.svc & {@env(hide)}},
		]
	}

	vols: {
		tmp: env.#Cache & {@env(docker-tmp), name: "docker-tmp"}
		varlib: env.#Cache & {@env(docker-varlib), name: "docker-varlib"}
	}
}
