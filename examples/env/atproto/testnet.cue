@experiment(aliasv2)
package atproto

import (
	"github.com/hofstadter-io/hof/examples/env/atproto/patches"
	"github.com/hofstadter-io/hof/lib/env/common/bases"
	"github.com/hofstadter-io/hof/lib/env/common/packs/databases"
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
					dir: env.#Dir & {
						path: "/work"
						source: env.#Container & {
							from: bases.debian
							steps: [env.Bash & {script: _relay}, env.Bash & {script: _pds}]
						}
					}
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
				from: builds.plc.ctr
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
		secret: env.#HostFile & {@env(), path: "./env/relay.secret.env"} // todo, we need secret version of this
		server: env.#Service & {
			@env()
			ports: [{port: 3000}]
			source: env.#Container & {
				from: builds.relay.ctr
				steps: [
					env.Envfile & {file: relay.config},
					env.Envfile & {file: relay.secret}, // todo, we need secret version of this
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
				from: builds.jetstream.ctr
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

	// @bluesky/pds or @blebbit/permissioned-pds
	pds: {
		config: env.#HostFile & {@env(), path: "./env/pds.env"}
		secret: env.#HostFile & {@env(), path: "./env/pds.secret.env"} // todo, we need secret version of this
		server: env.#Service & {
			@env()
			ports: [{port: 7002}]
			source: env.#Container & {
				from: _ | *builds.pds.ctr
				if _flags.blebbit {
					from: builds.ppds.ctr
				}
				steps: [
					env.Envfile & {file: pds.config},
					env.Envfile & {file: pds.secret}, // todo, we need secret version of this
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
		code: env.#Dir & {source: repos.blebbit}
		ctr: env.#DockerBuild & {source: code, dockerfile: "services/pds/Dockerfile"}
	}
	pds: {
		code: env.#Dir & {source: repos.atproto}
		ctr: env.#DockerBuild & {source: code, dockerfile: "services/pds/Dockerfile"}
	}
	plc: {
		// source
		code: env.#Dir & {source: repos.didplc}
		fixd: env.#Dir & {source: repos.didplc, patch: patches.plc}
		// images
		ctr: env.#DockerBuild & {source: fixd, dockerfile: "packages/server/Dockerfile"}
		dev: env.#Container & {
			from: ctr
			steps: [
				env.User & {name: "root"},
				env.Workdir & {path: "/app"},
				env.Envfile & {file: testnet.plc.config},
				env.BindService & {service: testnet.plc.postgres},
				env.Entrypoint & {args: ["sh"]},
				env.DefaultTerm & {args: ["sh"]},
			]
		}

		// example of adhoc work to figure out and apply a patch
		origCtr: env.#DockerBuild & {source: code, dockerfile: "packages/server/Dockerfile"}
		origDev: env.#Container & {
			from: origCtr
			steps: [
				env.User & {name: "root"},
				env.Workdir & {path: "/app"},
				env.Exec & {args: ["apk", "add", "--update", "patch", "git"]},
				env.File & {path: "/app/plc.diff", content: patches.plc},
				env.Envfile & {file: testnet.plc.config},
				env.BindService & {service: testnet.plc.postgres},
				env.Entrypoint & {args: ["sh"]},
				env.DefaultTerm & {args: ["sh"]},
			]
		}
	}
	relay: {
		code: env.#Dir & {source: repos.indigo}
		ctr: env.#DockerBuild & {source: code, dockerfile: "cmd/relay/Dockerfile"}
	}
	jetstream: {
		code: env.#Dir & {source: repos.jetstream}
		ctr: env.#DockerBuild & {source: code}
	}

}
