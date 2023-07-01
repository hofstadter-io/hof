package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "fmt"
	on:   _ | *["push"]
	env: HOF_TELEMETRY_DISABLED: "1"

	jobs: formatter: {
		concurrency: {
			group:                "${{ github.workflow }}-${{ github.ref_name }}"
			"cancel-in-progress": true
		}
		"runs-on":   "ubuntu-latest"
		environment: "hof mod testing"
		strategy: {
			"fail-fast": false
			matrix: formatter: common.Formatters
		}

		steps: [
			common.Steps.checkout,
			common.Steps.vars,
			common.Steps.buildx.qemu,
			common.Steps.buildx.setup.linux,
			common.Steps.buildx.formatters,
		]
	}
}
