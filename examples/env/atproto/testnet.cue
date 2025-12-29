@experiment(aliasv2)
package devex

import (
	"github.com/hofstadter-io/hof/lib/env/common/packs/databases"
	"github.com/hofstadter-io/hof/schemas/env"
)

testnet: {
	// give things consistent names
	[string]~(group,_): [string]~(subgroup,_): name: "\(group)-\(subgroup)"

	// @atproto PLC
	plc: {
		config: env.#HostFile & {@env(), path: "./env/plc.env"}
		server: env.#Service & {
			@env()
			ports: [{port: 3000}]
			source: env.#Container & {
				from: "blebbit/plc:latest"
				steps: [
					env.Envfile & {file: plc.config},
					env.BindService & {service: plc.postgres},
				]
			}
		}
		(databases.Postgres & {#name: "plc"}).#out
	}

	// @atproto Relay
	relay: {
		config: env.#HostFile & {@env(), path: "./env/relay.env"}
		server: env.#Service & {
			@env()
			ports: [{port: 3000}]
			source: env.#Container & {
				from: "blebbit/relay:latest"
				steps: [
					env.Envfile & {file: relay.config},
					env.Mount & {path: "/data", source: relay.data},
					env.BindService & {service: relay.postgres},
					env.BindService & {service: plc.server},
				]
			}
		}
		data: env.#Cache & {@env()}
		(databases.Postgres & {#name: "relay"}).#out
	}

	// @atproto Jetstream
	jetstream: {
		config: env.#HostFile & {@env(), path: "./env/jetstream.env"}
		server: env.#Service & {
			@env()
			ports: [{port: 7002}]
			source: env.#Container & {
				from: "blebbit/jetstream:latest"
				steps: [
					env.Envfile & {file: jetstream.config},
					env.Mount & {path: "/data", source: jetstream.data},
					env.BindService & {service: plc.server},
					env.BindService & {service: relay.server},
				]
			}
		}
		data: env.#Cache & {@env()}
	}

	// @blebbit Permissioning PDS
	pds: {
		config: env.#HostFile & {@env(), path: "./env/pds.env"}
		server: env.#Service & {
			@env()
			ports: [{port: 7002}]
			source: env.#Container & {
				from: "blebbit/jetstream:latest"
				steps: [
					env.Envfile & {file: pds.config},
					env.Mount & {path: "/app/data", source: pds.data},
					env.Mount & {path: "/app/blobs", source: pds.blobs},
					env.BindService & {service: pds.spicedb},
					env.BindService & {service: plc.server},
					env.BindService & {service: relay.server},
				]
			}
		}
		data: env.#Cache & {@env()}
		blobs: env.#Cache & {@env()}

		spicedb: env.#Service & {
			@env()
			ports: [{port: 8080}, {port: 9090}, {port: 50051}]
			args: ["serve", "--http-enabled"]
			source: env.#Container & {
				from: "authzed/spicedb:latest"
				envs: {
					SPICEDB_GRPC_PRESHARED_KEY: "testnet-spicedb"
					SPICEDB_DATASTORE_ENGINE:   "postgres"
					SPICEDB_DATASTORE_CONN_URI: "postgres://spicedb:spicedb@pds-pg-svc:5432/spicedb?sslmode=disable"
				}
				steps: [
					env.Exec & {args: ["migrate", "head"], useEntrypoint: true},
					env.BindService & {service: pds.postgres},
				]
			}
		}
		(databases.Postgres & {#name: "spicedb"}).#out
	}
}
