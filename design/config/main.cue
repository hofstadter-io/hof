package config

import (
	"github.com/hofstadter-io/hofmod-cuefig/schema"
)

#HofConfig: schema.#Config & {
  Name: "config"
  Entrypoint: ".hofcfg.cue"

  ConfigSchema: {
    models: {
      name: "string"
    }
    stores: {
      name: "string"
      type: "string"
    }
  }
}

#HofCredentials: schema.#Config & {
  Name: "secret"
  Entrypoint: ".hofshh.cue"

  ConfigSchema: {
    [Cred=string]: {
      [Key=string]: string
    }
  }
}
