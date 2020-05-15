package design

import (
	"github.com/hofstadter-io/hofmod-cuefig/schema"
)

#HofConfig: schema.#Config & {
  Name: "config"
  Entrypoint: "hofconfig.cue"

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
  Name: "creds"
  Entrypoint: "\(#CLI.ConfigDir)/credentials.cue"

  ConfigSchema: {
    [Cred=string]: {
      [Key=string]: string
    }
  }
}
