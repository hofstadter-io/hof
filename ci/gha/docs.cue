package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "docs"

	on: push: {
		paths: ["docs/**", "ci/gha/docs.cue"]
	}

	jobs: {
		docs: {
			"runs-on": "ubuntu-latest"
			// environment: "hof docs"

			steps: [
				common.Steps.go.setup & { #ver: "1.20.x" },
				common.Steps.go.cache,
				common.Steps.checkout,
				common.Steps.vars,
				common.Steps.go.deps,
				{
					name: "Build hof"
					run:  "go install ./cmd/hof"
				},
				{
					name: "Setup"
					run:  """
					cd docs
					make first 
					"""
				},
				{
					name: "Test"
					run:  """
					cd docs
					make test
					"""
				},
				{
					name: "Build"
					run:  """
					cd docs
					make build
					"""
				},

				{
					name: "Check"
					run:  """
					cd docs
					make run &
					make broken-link
					"""
				},
			]
		}
	}
}

