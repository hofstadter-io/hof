package containers

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

registry: {

	service: env.#Service & {
		@env(registry-svc)
		name: "registry"
		ports: [{name: "http", port: 5000}]
		source: container
	}

	container: env.#Container & {
		@env(registry-ctr)
		name: "registry"
		from: "registry:3"
		steps: [
			env.Mount & {path: "/var/lib/registry", source: data},
		]
	}

	data: env.#Cache & {
		@env(registry-data)
		name: "registry"
	}

	// https://github.com/distribution/distribution/blob/main/cmd/registry/config-example.yml
	// config: _
}
