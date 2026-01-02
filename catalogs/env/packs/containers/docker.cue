package containers

import "github.com/hofstadter-io/hof/schemas/env"

docker: {
	#ver: string | *"29"

	imgs: {
		_tags: [#ver, "cli", "dind", "dind-rootless"]
		for _, tag in _tags {
			(tag): env.#Container & {
				name: tag
				if tag == #ver {from: "docker:\(#ver)"}
				if tag != #ver {from: "docker:\(#ver)-\(tag)"}
			}
		}
	}

	cli: env.#File & {
		path:   "/usr/local/bin/docker"
		source: imgs.cli
	}

	daemon: {
		ctr: env.#Container & {
			from: imgs.dind
			steps: [
				env.Mount & {path: "/tmp", source: vols.tmp},
				env.Mount & {path: "/var/lib/docker", source: vols.lib},
				env.Expose & {port: 2375},
				env.Entrypoint & {args: [
					"dockerd",
					"--log-level=warn",
					"--host=tcp://0.0.0.0:2375",
					"--tls=false",
				]},
			]
		}

		svc: env.#Service & {
			name:   "global-dockerd"
			source: daemon.ctr
			ports: [{name: "docker", port: 2375}]

		}
	}

	steps: {
		bind: [
			env.Env & {DOCKER_HOST: "tcp://gbolal-dockerd:2375"},
			env.BindService & {service: daemon.svc},
		]
	}

	vols: {
		tmp: env.#Cache
		lib: env.#Cache
	}
}
