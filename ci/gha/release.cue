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
	env: HOF_TELEMETRY_DISABLED: "1"

	jobs: {
		goreleaser: {
			environment: "hof mod testing"
			"runs-on":   "ubuntu-latest"
			steps: [
				{
					name: "cancel if not our repository"
					run: """
						gh run cancel ${{ github.run_id }}	
						gh run watch  ${{ github.run_id }}	
						"""
					env: GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
					if: "github.repository != hofstadter-io/hof"
				},
				common.Steps.checkout,
				common.Steps.vars,
				common.Steps.buildx.qemu,
				common.Steps.buildx.setup.linux,
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
			"runs-on":   "ubuntu-latest"

			steps: [
				common.Steps.checkout,
				common.Steps.vars,
				common.Steps.buildx.qemu,
				common.Steps.buildx.setup.linux,
				common.Steps.docker.login,
				common.Steps.buildx.formatters & {with: push: true},
			]
		}
	}
}
