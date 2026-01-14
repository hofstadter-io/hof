package bases

import (
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

wolfi: {
	minimal: env.#Container & {
		@env(wolfi-min) // still need to do this, so that AtMade is true for viewing
		from: "cgr.dev/chainguard/wolfi-base"
		steps: [
			// globally consistent for us
			env.Workdir & {path: "/work"},
			// stuff we don't want in the final image, but need around when building/running
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "wolfi-var-log"}},
			env.Mount & {path: "/var/cache", source: env.#Cache & {name: "wolfi-var-cache"}},

			// update just once at the beginning
			utils.apk.update,
			// utils.apk.upgrade,

			// the minimal essenitals
			utils.apk.install & {#pkgs: ["ca-certificates", "wget", "curl", "bash"]},
		]
	}
}

chaingaurd: {
	// free
	apko: {} // build images from apk w/o dockerfile
	"aspnet-runtime": {}
	bash: {}
	busybox: {}
	cosign: {}
	crane: {}
	curl: {}
	dfc: {}
	dive: {}
	"docker-dind": {}
	"dotnet-runtime": {}
	"dotnet-sdk": {}
	"gcc-glibc": {}
	git: {}
	"glibc-dynamic": {}
	go: {}
	gradle: {}
	grype: {} // vuln scanning
	haproxy: {}
	helm: {}
	hugo: {}
	jdk: {}
	jre: {}
	k3s: {}
	"k3s-static": {}
	ko: {}
	kubectl: {}
	laravel: {}
	malconent: {} // secops
	mariadb: {}
	maven: {}
	mongodb: {}
	nginx: {}
	node: {}
	openscap: {}
	"opentelemetry-collector-contrib": {}
	php: {}
	postgres: {}
	python: {}
	pytorch: {}
	redis: {}
	ruby: {}
	rust: {}
	static: {}
	valkey: {}
	"wait-for-it": {}
	"wolfi-base": {}
	wordpress: {}
}