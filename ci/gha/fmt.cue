package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "fmt"

	_on:   ["push", "pull_request"]
	_paths: ["lib/fmt/**", "formatters/**", "ci/gha/fmt.cue", ".github/workflows/fmt.yml"]
	on: { for evt in _on { (evt): paths: _paths } }
	on: workflow_dispatch: {}
	env: HOF_TELEMETRY_DISABLED: "1"


	jobs: formatter: {
		"runs-on": "ubuntu-latest"
		environment: "hof mod testing"
		strategy: {
			"fail-fast": false
			matrix: formatter: common.Formatters
		}

		steps: [
			common.Steps.checkout,
			common.Steps.vars,
			common.Steps.buildx.qemu,
			common.Steps.buildx.setup,
			common.Steps.buildx.formatters,
		]
	}
}
