package workflows

import "github.com/hofstadter-io/ghacue"

_bases: ["debian"]
_formatters: [
	"prettier",
	"black",
]

ghacue.#Workflow & {
	name: "fmt"

	_on:   ["push", "pull_request", "workflow_dispatch"]
	_paths: ["formatters/**",".github/workflows/fmt.*"]
	on: { for evt in _on { (evt): paths: _paths } }

	jobs: {
		for fmtr in _formatters {
			for base in _bases {
				"\(fmtr)-\(base)": {
					"runs-on": "ubuntu-latest"
					// environment: "hof mod testing"
					steps: [{
						name: "Checkout code"
						uses: "actions/checkout@v3"
					},{
						name: "Set up QEMU"
						uses: "docker/setup-buildx-action@v2"
					},{
						name: "Set up Docker BuildX"
						uses: "docker/setup-buildx-action@v2"
			// 		},{
			// 			name: "Login to DockerHub"
			//			if: "github.event_name != 'pull_request'"
			// 			uses: "docker/login-action@v2"
			// 			with: {
			// 				username: "${{ secrets.DOCKERHUB_USERNAME }}"
			// 				password: "${{ secrets.DOCKERHUB_TOKEN }}"
			// 			}
					},{
						name: "Build Image"
						uses: "docker/build-push-action@v3"
						with: {
							context: "formatters/tools/\(fmtr)"
							file: "\(context)/Dockerfile.\(base)"
							// push: true
							platforms: "linux/amd64,linux/arm64"
							// TODO, figure out tag
							tags: """
							"hofstadter-io/fmt-\(fmtr):dirty-\(base)"
							"""
						}
					}]
				}
			}
		}
	}
}
