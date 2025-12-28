@experiment(aliasv2)
package devex

import "github.com/hofstadter-io/hof/schemas/env"

_flags: {
	tier3: {
		repo:    string | *"https://github.com/prasad-moru/3tier_app" @tag(tier3repo)
		appPort: int | *3000                                          @tag(tier3_appPort)
		apiPort: int | *3000                                          @tag(tier3_apiPort)
	}
}

tier3: {
	[string]~(comp,_): name: "tier3-\(comp)"

	// source code
	repo: env.#GitRepo & {
		@env()
		url: _flags.tier3.repo
	}

	// frontend app
	app: env.#Service & {
		@env()
		ports: [{port: _flags.tier3.appPort}]
		source: tier3.appCtr
	}

	// backend api
	api: env.#Service & {
		@env()
		ports: [{port: _flags.tier3.apiPort}]
		source: tier3.apiCtr
	}

	// database server & volume
	db: env.#Service & {
		@env()
		ports: [{port: 3306}]
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
		@env()
		from: "node:24"
		steps: [
      (_nodePrep & { dir: "frontend" }).steps,
			env.Exec & {args: ["npm", "run", "build"]},
			env.Args & {args: ["npx", "serve", "-s", "build"]},
			env.Expose & {port: _flags.tier3.appPort, name: "http"},
			env.BindService & {service: tier3.api},
		]
	}
	apiCtr: env.#Container & {
		@env()
		from: "node:24"
		steps: [
      (_nodePrep & { dir: "backend" }).steps,
			env.Args & {args: ["npm", "start"]},
			env.Expose & {port: _flags.tier3.apiPort, name: "http"},
			env.BindService & {service: tier3.db},
		]
	}

  // sharing is caring (function pattern)
  _nodePrep: {
    dir: string
    _dir: env.#Dir & { path: dir, source: tier3.repo}
    steps: [
			env.Workdir & {path: "/app"},
			env.Dir & {path: "/app", source: _dir, include: ["package*.json"]},
			env.Exec & {args: ["npm", "install"]},
			env.Dir & {path: "/app", source: _dir},
    ]
  }

  // dev container
	appTest: env.#Container & {
		@env()
		from: "node:24"
		steps: [
      // (_nodePrep & { dir: "frontend" }).steps,
      // (_nodePrep & { dir: "backend" }).steps,
			env.BindService & {service: tier3.db},
			env.BindService & {service: tier3.api},
			env.BindService & {service: tier3.app},
      env.Entrypoint & {args: ["sh"], }
    ]
  }
}
