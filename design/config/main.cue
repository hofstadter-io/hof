package config

import (
	hof "github.com/hofstadter-io/hof/schema"

	"github.com/hofstadter-io/hofmod-cuefig/schema"
)

#HofConfig: schema.#Config & {
  Name: "config"
  Entrypoint: ".hofcfg.cue"

  ConfigSchema: {
    Modelset: hof.#Modelsets
    Datastores: hof.#Datastores
  }
}

#HofCredentials: schema.#Config & {
  Name: "secret"
  Entrypoint: ".hofshh.cue"

  ConfigSchema: {
    [Group=string]: {
      [Cred=string]: {
        [Key=string]: string
      }
    }
  }
}
