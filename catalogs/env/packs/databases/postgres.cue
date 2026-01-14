package databases

import "github.com/hofstadter-io/hof/schemas/env"

Postgres: {
	#name: string
	#port: int | *5432
	#ver: string | *"18"
	#temp: bool | *false

	volume: env.#Cache & {name: string | *"\(#name)-pg-data"}

	container: env.#Container & {
		name: string | *"\(#name)-pg"
		from: "postgres:\(#ver)"
		envs: {
			POSTGRES_DB:       #name
			POSTGRES_PORT:     "\(#port)"
			POSTGRES_USER:     #name
			POSTGRES_PASSWORD: #name
		}
		steps: [
			if !#temp { env.Mount & {path: "/var/lib/postgresql", source: volume} },
			if  #temp { env.Temp & {path: "/var/lib/postgresql"} },
		]
	}

	service: env.#Service & {
		hostname: string | *"\(#name)-pg"
		ports: [{port: #port}]
		source: container
	}

	#out: {
		postgresVolume: volume
		postgres: service
	}
}
