package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/.github/workflows/common"
)

_formatters: [
	"prettier",
	"csharpier",
	"black",
]

ghacue.#Workflow & {
	name: "release"

	on: push: {
		"branches-ignore": ["*", "!_dev"]
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
			environment: "hof mod testing"
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
				common.Steps.docker.login,
				common.Steps.docker."fmtr-buildx" & { with: push: true },
			]
		}
	}
}
