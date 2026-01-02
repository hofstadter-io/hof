package tool

import (
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

docker: {
	repo: env.Exec & {
		args: ["sh", "-c", _script]

		_script: """
			set -eou pipefail

			install -m 0755 -d /etc/apt/keyrings
			curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
			chmod a+r /etc/apt/keyrings/docker.asc


			cat <<EOF > /etc/apt/sources.list.d/docker.sources
			Types: deb
			URIs: https://download.docker.com/linux/debian
			Suites: trixie
			Components: stable
			Signed-By: /etc/apt/keyrings/docker.asc
			EOF

			apt-get update -y
			"""
	}

	cli: [
		docker.repo,
		utils.apt.install & {#pkgs: ["docker-ce-cli", "docker-buildx-plugin", "docker-compose-plugin"]},
	]

}
