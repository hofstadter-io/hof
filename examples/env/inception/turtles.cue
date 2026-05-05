@experiment(aliasv2)
package inception

import (
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/catalogs/env/packs"
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

// local pack with current os & arch set as defaults
_packs: packs & {flags: root.flags}

turtles: {
	// any container
	dev: env.#Container & {
		@env(turtles-dev)
		from: bases.debian13.default

		steps: [
			// use zsh
			_packs.tool.zsh.customize,

			// add hof/veg
			// _veg.hof.File["linux-arm64"],

			// add a bunch of tools
			_packs.containers.docker.cli.install,
			_packs.containers.dagger.cli.install,
			_packs.tool.k8s.kubectl,
			_packs.tool.k8s.helm,
			_packs.tool.k8s.crane,
			// _packs.tool.k8s.kind.binary,

			// add the socket for inception
			env.UnixSocket & {path: "/var/run/docker.sock", source: turtles.socket},
		]
	}

	// get a socket from the host, inception will be painfully slow otherwise
	socket: env.#HostSocket & {
		@env(turtles-sock)
		name: "turtles-sock"
		path: flags.socket
	}

}

flags: {

	goos: string @tag(goos,var=os)
	arch: string @tag(arch,var=arch)

	socket: string | *"unix:///var/run/docker.sock" @tag(socket)
}
