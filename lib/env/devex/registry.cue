package devex

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

registry: {
  container: env.Container & {
    @env()
    name: "registry"
    from: "registry:3"
  }

  service: env.Service & {
    @env()
    name: "registry"
    port: 5000
    image: container
    volumes: [volume] // todo, config?
  }

  volume: env.Volume & {
    @env()
    name: "registry"
    type: "cache"
  }

}
