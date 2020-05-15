package hofconfig

import (
	"github.com/hofstadter-io/hof/schema"
)

modelsets: schema.#Modelsets & {
  test: {
    entry: "models"

    stores: {
      test: "test"
      dev:  "dev"
    }
  }
  other: {
    entry: "others" @entry(hof=geb)

    stores: {
      test: "test" @hof("geb")
      dev:  "dev"
    }
  } @hof(geb=weird) @geb()
} @model(hof=goo)

stores: schema.#Stores & {
  test: {
    id: "test"
    type: "psql"
    version: "12.2"
  }
  dev: {
    id: "local-dev"
    type: "psql"
    version: "12.2"
  }
}
