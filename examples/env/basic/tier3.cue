@experiment(aliasv2)
package basic

import "github.com/hofstadter-io/hof/schemas/env"

_flags: {
	tier3: {
		repo:    string | *"https://github.com/prasad-moru/3tier_app" @tag(tier3_repo)
		appPort: int | *3000                                          @tag(tier3_appPort)
		apiPort: int | *3001                                          @tag(tier3_apiPort)
		dbPort:  int | *3306                                          @tag(tier3_dbPort)
	}
}

tier3: {
	[string]~(comp,_): name: "tier3-\(comp)"

	// source code
	repo: env.#GitRepo & {
		@env(tier3-repo)
		url: _flags.tier3.repo
	}

	// frontend app
	app: env.#Service & {
		@env(tier3-app)
		ports: [{port: 3000, frontend: _flags.tier3.appPort}]
		source: tier3.appCtr
	}

	// backend api
	api: env.#Service & {
		@env(tier3-api)
		ports: [{port: 3000, frontend: _flags.tier3.apiPort}]
		source: tier3.apiCtr
	}

	// database server & volume
	db: env.#Service & {
		@env(tier3-db)
		ports: [{port: 3306, frontend: _flags.tier3.dbPort}]
		source: env.#Container & {
			name: "mysql"
			from: "mysql:8.0"
			envs: {
				MYSQL_DATABASE:      "appdb"
				MYSQL_PASSWORD:      "pass123"
				MYSQL_ROOT_PASSWORD: "pass123"
			}
		}
	}

	// containers, built and prepared
	appCtr: env.#Container & {
		from: "node:24"
		steps: [
			(_nodePrep & {dir: "frontend"}).steps,
			env.Exec & {args: ["npm", "run", "build"]},
			env.DefaultArgs & {args: ["npx", "serve", "-s", "build"]},
			env.Expose & {port: 3000, name: "http"},
			env.BindService & {service: tier3.api},
		]
	}
	apiCtr: env.#Container & {
		from: "node:24"
		steps: [
			(_nodePrep & {dir: "backend"}).steps,
			env.DefaultArgs & {args: ["npm", "start"]},
			env.Expose & {port: 3000, name: "http"},
			env.BindService & {service: tier3.db},
		]
	}

	// sharing is caring (function pattern)
	_nodePrep: {
		dir: string
		_dir: env.#Dir & {path: dir, sources: [tier3.repo]}
		steps: [
			env.Workdir & {path: "/app/\(dir)"},
			env.Dir & {path: "/app", source: _dir},
			env.Exec & {args: ["npm", "install"]},
		]
	}

	// dev container
	play: env.#Container & {
		@env(tier3-play)
		from: "node:24"
		steps: [
			// (_nodePrep & { dir: "frontend" }).steps,
			// (_nodePrep & { dir: "backend" }).steps,
			env.BindService & {service: tier3.db},
			env.BindService & {service: tier3.api},
			env.BindService & {service: tier3.app},
			env.Entrypoint & {args: ["sh"]},
		]
	}
}
