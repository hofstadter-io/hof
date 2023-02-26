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
					make test
					"""
				},
				{
					name: "Build"
					run:  """
					make build
					"""
				},

				{
					name: "Check"
					run:  """
					make run &
					make broken-link
					"""
				},
			]
		}
	}
}

