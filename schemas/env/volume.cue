package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

Volume: {
	schemas.Hof
	#hof: env: {
		root: true
		kind: "volume"
	}

  // host is only available in docker / run
  name: string
  type: "cache" | "host"
}