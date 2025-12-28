@experiment(aliasv2)
package devex

import "github.com/hofstadter-io/hof/schemas/env"

#pg: {
	#name: string
	#out: {
		postgresVolume~V: env.#Cache & {@env()}
		postgres: env.#Service & {
			@env()
			ports: [{port: 5432}]
			source: env.#Container & {
				from: "postgres:16"
				envs: {
					POSTGRES_DB:       #name
					POSTGRES_USER:     #name
					POSTGRES_PASSWORD: #name
				}
        steps: [
					env.Mount & {path: "/var/lib/postgresql/data", source: V},
        ]
			}
		}
	}
}