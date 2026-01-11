package containers

import "github.com/hofstadter-io/hof/catalogs/env/utils"

buildah: {
	// hack for now to maintain consistency in pack.<tool>.cli.install
	cli: install: [
		utils.apt.install & {#pkgs: ["buildah"]},
	]
}
