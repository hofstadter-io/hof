@experiment(aliasv2)
package atproto

import (
	"github.com/hofstadter-io/hof/examples/env/atproto/patches"
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/catalogs/env/packs/databases"
	"github.com/hofstadter-io/hof/schemas/env"
)

_flags: {
	blebbit: bool | *false @tag(blebbit,short=blebbit)
}

cmd: {
	[string]~(k1,_): env.#Cmd & {
		// @env(), name: k1
		tasks: [string]~(k2,_): {
			steps: [...[...{name: "\(k1).\(k2)"}]]
		}
	}

	init: tasks: {
		secrets: {
			steps: [[
				env.#ExportDir & {
					@env()
					path: "./env"
					sources: [env.#Dir & {
						path: "/work"
						sources: [
							env.#Container & {
								from: bases.debian.default
								steps: [env.Bash & {script: _relay}, env.Bash & {script: _pds}]
							},
						]
					}]
				},
			]]

			_relay: """
				(
				  echo RELAY_ADMIN_PASSWORD=$(openssl rand -hex 32 | tr -d '\n')
				) > relay.secret.env
				"""
			_pds: """
				(
				  echo # Private keys - these are each expected to be 64 char hex strings (1024 bit)
				  echo PDS_REPO_SIGNING_KEY_K256_PRIVATE_KEY_HEX=$(openssl rand -hex 64 | tr -d '\n')
				  echo PDS_PLC_ROTATION_KEY_K256_PRIVATE_KEY_HEX=$(openssl rand -hex 64 | tr -d '\n')
				  echo # Secrets - update to secure high-entropy strings
				  echo PDS_DPOP_SECRET=$(openssl rand -hex 32 | tr -d '\n')
				  echo PDS_JWT_SECRET=$(openssl rand -hex 32 | tr -d '\n')
				  echo PDS_ADMIN_PASSWORD=$(openssl rand -hex 32 | tr -d '\n')
				  echo PDS_SPICEDB_TOKEN=$(openssl rand -hex 32 | tr -d '\n')
				) > pds.secret.env
				"""
		}
	}
}

// naming things is not so hard

_testnet: [string]~(S,_): {
	service?: {@env(), name: "\(S)"}
	config?: {@env(), name: "\(S)-cfg"}
}
_testnet: [=~"(relay|pds)"]~(S2,_): {
	secret?: {@env(), name: "\(S2)-shh"}
}
_testnet: [!~"(jetstream)"]~(S3,_): {
	postgres?: {@env(), name: "\(S3)-pg"}
}

testnet: _testnet & {
	// @atproto PLC
	plc: {
		config: env.#HostFile & {path: "./env/plc.env"}
		service: env.#Service & {
			ports: [{port: 3000}]
			source: env.#Container & {
				from: builds.plc.ctr
				steps: [
					env.EnvFile & {file: plc.config},
					env.BindService & {service: plc.postgres},
				]
			}
		}
		(databases.Postgres & {#name: "plc"}).#out
	}

	// @atproto Relay
	relay: {
		config: env.#HostFile & {path: "./env/relay.env"}
		secret: env.#HostFile & {path: "./env/relay.secret.env"}
		service: env.#Service & {
			ports: [{port: 3000}]
			source: env.#Container & {
				from: builds.relay.ctr
				steps: [
					env.EnvFile & {file: relay.config},
					env.EnvFile & {file: relay.secret}, // todo, we need secret version of this
					env.Mount & {path: "/data", source: relay.data},
					env.BindService & {service: relay.postgres},
					env.BindService & {service: plc.service},
				]
			}
		}
		data: env.#Cache & {name: "relay-data"}
		(databases.Postgres & {#name: "relay"}).#out
	}

	// @atproto Jetstream
	jetstream: {
		config: env.#HostFile & {path: "./env/jetstream.env"}
		service: env.#Service & {
			ports: [{port: 3000}]
			source: env.#Container & {
				from: builds.jetstream.ctr
				steps: [
					env.EnvFile & {file: jetstream.config},
					env.Mount & {path: "/data", source: jetstream.data},
					env.BindService & {service: plc.service},
					env.BindService & {service: relay.service},
				]
			}
		}
		data: env.#Cache & {name: "jetstream-data"}
	}

	// @bluesky/pds or @blebbit/permissioned-pds
	pds: {
		config: env.#HostFile & {path: "./env/pds.env"}
		// TODO, #Secret (make and then provide to #SecretEnvfile)
		secret: env.#HostFile & {path: "./env/pds.secret.env"}
		service: env.#Service & {
			ports: [{port: 3000}]
			source: env.#Container & {
				from: _ | *builds.pds.ctr
				if _flags.blebbit {
					from: builds.ppds.ctr
				}
				steps: [
					env.EnvFile & {file: pds.config},
					env.EnvFile & {file: pds.secret}, // todo, we need secret version of this
					env.Mount & {path: "/app/data", source: pds.data},
					env.Mount & {path: "/app/blobs", source: pds.blobs},
					env.BindService & {service: pds.spicedb},
					env.BindService & {service: plc.service},
					env.BindService & {service: relay.service},
				]
			}
		}
		data: env.#Cache & {name: "pds-data"}
		blobs: env.#Cache & {name: "pds-blobs"}

		spicedb: env.#Service & {
			name: "pds-spicedb"
			ports: [{port: 8080}, {port: 9090}, {port: 50051}]
			args: ["serve", "--http-enabled"]
			source: env.#Container & {
				from: "authzed/spicedb:latest"
				envs: {
					SPICEDB_GRPC_PRESHARED_KEY: "testnet-spicedb"
					SPICEDB_DATASTORE_ENGINE:   "postgres"
					SPICEDB_DATASTORE_CONN_URI: "postgres://spicedb:spicedb@pds-pg:5432/spicedb?sslmode=disable"
				}
				steps: [
					env.Exec & {args: ["migrate", "head"], useEntrypoint: true},
					env.BindService & {service: pds.postgres},
				]
			}
		}
		(databases.Postgres & {#name: "pds-spicedb"}).#out
	}
}

builds: {
	// give things consistent names
	[string]~(group,_): [string]~(subgroup,_): {@env(), name: "\(group)-\(subgroup)"}

	repos: {
		blebbit: env.#GitRepo & {url: "https://github.com/blebbit/atproto"}
		atproto: env.#GitRepo & {url: "https://github.com/bluesky-social/atproto"}
		didplc: env.#GitRepo & {url: "https://github.com/did-method-plc/did-method-plc"}
		indigo: env.#GitRepo & {url: "https://github.com/bluesky-social/indigo"}
		jetstream: env.#GitRepo & {url: "https://github.com/bluesky-social/jetstream"}
	}
	ppds: {
		code: env.#Dir & {sources: [repos.blebbit]}
		ctr: env.#DockerBuild & {source: code, dockerfile: "services/pds/Dockerfile"}
	}
	pds: {
		code: env.#Dir & {sources: [repos.atproto]}
		ctr: env.#DockerBuild & {source: code, dockerfile: "services/pds/Dockerfile"}
	}
	plc: {
		code: env.#Dir & {sources: [repos.didplc]}
		// patch
		fixd: env.#Dir & {sources: [repos.didplc], patch: patches.plc}
		ctr: env.#DockerBuild & {source: fixd, dockerfile: "packages/server/Dockerfile"}
	}
	relay: {
		code: env.#Dir & {sources: [repos.indigo]}
		// patch
		fixd: env.#Dir & {sources: [repos.indigo], patch: patches.relay}
		ctr: env.#DockerBuild & {source: fixd, dockerfile: "cmd/relay/Dockerfile"}
	}
	jetstream: {
		code: env.#Dir & {sources: [repos.jetstream]}
		ctr: env.#DockerBuild & {source: code}
	}

	// hack: {
	// 	ctr: env.#Container

	// 	// images
	// 	dev: env.#Container & {
	// 		from: ctr
	// 		steps: [
	// 			env.User & {name: "root"},
	// 			env.Workdir & {path: "/app"},
	// 			env.Envfile & {file: testnet.plc.config},
	// 			env.BindService & {service: testnet.plc.postgres},
	// 			env.Entrypoint & {args: ["sh"]},
	// 			env.DefaultTerm & {args: ["sh"]},
	// 		]
	// 	}
	// }
}
