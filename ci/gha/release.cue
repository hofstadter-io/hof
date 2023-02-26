package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "release (hof)"

	on: push: {
		tags: ["v*"]
	}

	jobs: {
		goreleaser: {
			environment: "hof mod testing"
			"runs-on": "ubuntu-latest"
			steps: [
				common.Steps.checkout,
				common.Steps.vars,
				common.Steps.docker.qemu,
				common.Steps.docker.setup,
				common.Steps.docker.login,
				common.Steps.go.setup,
				common.Steps.go.deps,
				common.Steps.go.releaser,
			]
		}
		formatter: {
			strategy: {
				"fail-fast": false
				matrix: formatter: common.Formatters
			}
			environment: "hof mod testing"
			"runs-on": "ubuntu-latest"

			steps: [
				common.Steps.checkout,
				common.Steps.vars,
				common.Steps.docker.qemu,
				common.Steps.docker.setup,
				common.Steps.docker.login,
				common.Steps.docker.formatters & { with: push: true },
			]
		}
	}
}
