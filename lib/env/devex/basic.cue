package devex

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

veg: env.Container & {
    @env()
    from: "debian:13-slim"
}