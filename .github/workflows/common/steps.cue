package common

import "strings"

Steps: {

	checkout: {
		name: "Checkout code"
		uses: "actions/checkout@v3"
	}

	vars: {
		name: "Setup Vars"
		id: "vars"
		run: """
		SHA=${GITHUB_SHA::8}
		TAG=$(git tag --points-at HEAD)
		if [ -z $TAG ]; then
			TAG=${SHA}
		fi
		echo "HOF_SHA=${SHA}" >> $GITHUB_ENV
		echo "HOF_TAG=${TAG}" >> $GITHUB_ENV
		"""
	}

	go: {
		setup: {
			name: "Install Go"
			uses: "actions/setup-go@v3"
			with: "go-version": "${{ matrix.go-version }}"
		}
		cache: {
			uses: "actions/cache@v3"
			with: {
				path: #"""
					~/go/pkg/mod
					~/.cache/go-build
					~/Library/Caches/go-build
					~\AppData\Local\go-build
					"""#
				key:            "${{ runner.os }}-go-${{ matrix.go-version }}-${{ hashFiles('**/go.sum') }}"
				"restore-keys": "${{ runner.os }}-go-${{ matrix.go-version }}-"
			}
		}
		deps: {
			name: "Fetch Go deps"
			run:  "go mod download"
		}
	}

	docker: {
		qemu: {
			name: "Set up QEMU"
			uses: "docker/setup-buildx-action@v2"
		}

		setup: {
			name: "Set up Docker BuildX"
			uses: "docker/setup-buildx-action@v2"
		}

		login: {
			name: "Login to Docker Hub"
			uses: "docker/login-action@v2"
			with: {
				username: "${{ secrets.HOF_DOCKER_USER }}"
				password: "${{ secrets.HOF_DOCKER_TOKEN }}"
			}
		}

		"fmtr-buildx": {
			name: "Build Image"
			uses: "docker/build-push-action@v3"
			with: {
				context: "formatters/tools/${{ matrix.formatter }}"
				file: "\(context)/Dockerfile.debian"
				platforms: "linux/amd64,linux/arm64"
				tags: strings.Join([
					"hofstadter-io/fmt-${{ matrix.formatter }}:${{ env.HOF_SHA }}",
					"hofstadter-io/fmt-${{ matrix.formatter }}:${{ env.HOF_TAG }}",
				], ",")
			}
		}
	}
}

