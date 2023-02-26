package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "release (docs)"

	on: push: {
		tags: ["docs-*"]
	}

	jobs: {
		docs: {
			"runs-on": "ubuntu-latest"
			// environment: "hof docs"

			steps: [
				common.Steps.checkout,
				common.Steps.vars,
				{
					name: "Setup"
					run:  "go install ./cmd/hof"
				},

			]
		}
	}
}

