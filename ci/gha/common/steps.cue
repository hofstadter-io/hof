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
		echo "HOF_FMT_VERSION=${TAG}" >> $GITHUB_ENV
		if [ -z $TAG ]; then
			TAG=${SHA}
		fi
		echo "HOF_SHA=${SHA}" >> $GITHUB_ENV
		echo "HOF_TAG=${TAG}" >> $GITHUB_ENV
		"""
	}

	go: {
		setup: {
			#ver: string | *(string & Versions.go)
			name: "Install Go"
			uses: "actions/setup-go@v3"
			with: "go-version": #ver
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
		releaser: {
			name: "Run GoReleaser"
			uses: "goreleaser/goreleaser-action@v4"
			with: {
				// either 'goreleaser' (default) or 'goreleaser-pro'
				distribution: "goreleaser"
				version: "latest"
				workdir: "cmd/hof"
				args: "release --clean -f goreleaser.yml -p 1"
			}
			env: {
				GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
			}
		}
	}

	docker: {
		qemu: {
			name: "Set up QEMU"
			uses: "docker/setup-qemu-action@v2"
			with: {
				platforms: "arm64"
			}
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

		formatters: {
			name: "Build Image"
			uses: "docker/build-push-action@v3"
			with: {
				context: "formatters/tools/${{ matrix.formatter }}"
				file: "\(context)/Dockerfile.debian"
				platforms: "linux/amd64,linux/arm64"
				tags: strings.Join([
					"hofstadter/fmt-${{ matrix.formatter }}:${{ env.HOF_SHA }}",
					"hofstadter/fmt-${{ matrix.formatter }}:${{ env.HOF_TAG }}",
				], ",")
			}
		}

		compat: {
			name: "Test Compatibility"
			run: """
				docker version
				go run test/docker/main.go
				"""

		}
	}
}

