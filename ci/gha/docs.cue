package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "docs"

	on: {
		for p in ["push"] {
			(p): {
				paths: [
					"docs/**",
					"ci/gha/docs.cue",
					"design/**",
					"schema/**",
					"cmd/**",
				]
			}
		}
	}
	env: HOF_TELEMETRY_DISABLED: "1"

	jobs: {
		docs: {
			"runs-on": "ubuntu-latest"

			steps: [
				// general setup
				common.Steps.cue.install,
				common.Steps.go.setup & {#ver: "1.20.x"},
				common.Steps.go.cache,
				common.Steps.checkout,
				common.Steps.vars,
				common.Steps.go.deps,
				common.Steps.hof.install,

				// dev build site & test
				common.Steps.docs.setup,
				{
					name: "Test"
					run: """
						cd docs
						make gen
						make test
						make run &
						make broken-link
						"""
				},
			]
		}
	}
}
