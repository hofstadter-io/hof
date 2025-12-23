package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

Service: {
	schemas.Hof
	#hof: env: {
		root: true
		kind: "service"
	}
}