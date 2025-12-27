package devex

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

registry: {

	service: env.#Service & {
		@env()
		name: "registry"
		ports: [{
      name: "http"
			port: 5000
		}]

		source: env.#Container & {
			@env()
			name: "registry"
			from: "registry:3"
		}

		// image: container
		// volumes: [volume] // todo, config?
	}

	volume: env.#Cache & {
		@env()
		name: "registry"
	}

}
