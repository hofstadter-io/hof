package docs

import (
	"github.com/hofstadter-io/supacode/gen"
	"github.com/hofstadter-io/supacode/schema"
)

// Generator definition
Generator: gen.Generator & {
	Name:   "docs"
	Outdir: "./"
	App:    schema.App & {
		name:   "docs"
		module: "github.com/hofstadter-io/hof"

		search: enabled:   false
		auth: enabled:     false
		database: enabled: false
	}

	Datamodel: {}

}

Workflows: {
	dev: {
		@flow(dev)

		first: {
			@task(os.Exec)
			cmd: ["hof", "gen", "@supacode"]
			exitcode: _
		}
		hof: {
			dep: first

			@task(os.Watch)
			globs: [
				"*.cue",
				"../cue.mod/pkg/github.com/hofstadter-io/supacode/**/*.*",
			]
			handler: {
				event?: _
				compile: {
					@task(os.Exec)
					cmd: ["hof", "gen", "@supacode"]
					exitcode: _
				}
				now: {
					dep: compile.exitcode
					n:   string @task(gen.Now)
					s:   "\(n) (\(dep))"
				}
				alert: {
					@task(os.Stdout)
					dep:  now.s
					text: "hof regen: \(now.s)\n"
				}
			}
		}

		npm: {
			dep: first
			@task(os.Exec)
			cmd: ["npm", "run", "dev"]
			exitcode: _
		}
	}
}
