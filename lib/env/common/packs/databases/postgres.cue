@experiment(aliasv2)
package databases

import "github.com/hofstadter-io/hof/schemas/env"

Postgres: {
	#name: string
	#port: int | *5432

	volume: env.#Cache & {name: string | *"\(#name)-pg-data"}

	container: env.#Container & {
		name: string | *"\(#name)-pg"
		from: "postgres:16"
		envs: {
			POSTGRES_DB:       #name
			POSTGRES_PORT:     "\(#port)"
			POSTGRES_USER:     #name
			POSTGRES_PASSWORD: #name
		}
		steps: [
			env.Mount & {path: "/var/lib/postgresql/data", source: volume},
		]
	}

	service: env.#Service & {
		name: string | *"\(#name)-pg-svc"
		ports: [{port: #port}]
		source: container
	}

	#out: {
		postgresVolume: volume & {@env()}
		postgres: service & {@env()}
	}
}
