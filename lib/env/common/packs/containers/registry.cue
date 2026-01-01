package containers

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

registry: {

	service: env.#Service & {
		@env()
		name: "registry"
		ports: [{name: "http", port: 5000}]
		source: container
	}

	container: env.#Container & {
		@env()
		name: "registry"
		from: "registry:3"
		steps: [
			env.Mount & {path: "/var/lib/registry", source: data},
		]
	}

	data: env.#Cache & {
		@env()
		name: "registry"
	}

}
