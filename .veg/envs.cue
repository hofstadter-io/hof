@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

_reg: "host.docker.internal:5000"
_ver: "latest"
// _reg: "ghcr.io/hofstadter-io"
// _ver: "v0.7.0-alpha.2"

environs: {
  for short in ["dev", "hof"] {
    (short): {
			@agentic(environ)
      name: "veg-\(short)"
      description: "our veg-\(short) image for self & agents"
      spec: {
        from: "\(_reg)/veg-\(short):\(_ver)"
      }
    }
  }

  agent: {
		@agentic(environ)
    name: "veg-agent"
    description: "optimized image for working on / with veg"
    specValue: root.ctr.agent
  }
  
}

ctr: {

	agent: env.#Container & {
		@env()
		#hof: {
			id: "veg-agent"
			metadata: {
				name:        id
				description: "runtime veg-dev, with sockets, but no secrets"
			}
		}
		name: #hof.metadata.name

		from: root.ctr.dev
		steps: [
			// add hof late (not in dev), because it changes frequently
			root.hof.File.linux,

			// config / env stuff
			// _packs.tool.k8s.kind.config,
			// going to switch to k3d / k3s

			// add the socket for inception
			env.UnixSocket & {path: "/var/run/docker.sock", source: root.host.docker.socket},

			// // bind lsp servers, started on demand
			// env.BindService & {service: lang.go.lsp},
			// env.BindService & {service: lang.cue.lsp},
			// env.BindService & {service: lang.node.lsp},
			// env.BindService & {service: lang.python.lsp},

		]
	}

}
