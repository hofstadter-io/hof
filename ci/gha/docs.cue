package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "docs"

	on: _ | *["push"]
	env: HOF_TELEMETRY_DISABLED: "1"

	jobs: {
		docs: {
			concurrency: {
				group:                "${{ github.workflow }}-${{ github.ref_name }}"
				"cancel-in-progress": true
			}
			"runs-on": "ubuntu-latest"

			steps: [
				// general setup
				common.Steps.cue.install,
				common.Steps.go.setup,
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
						hof fmt start prettier@v0.6.8
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
