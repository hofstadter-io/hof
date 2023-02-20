package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "fmt"

	_on:   ["push", "pull_request"]
	_paths: ["lib/fmt/**", "formatters/**","ci/gha/fmt.*"]
	on: { for evt in _on { (evt): paths: _paths } }
	on: workflow_dispatch: {}


	jobs: formatter: {
		"runs-on": "ubuntu-latest"
		strategy: {
			"fail-fast": false
			matrix: formatter: common.Formatters
		}

		steps: [
			common.Steps.checkout,
			common.Steps.vars,
			common.Steps.docker.qemu,
			common.Steps.docker.setup,
			common.Steps.docker.formatters,
		]
	}
}
