package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

#CI: {
	schemas.Hof

	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "ci"
	}

	// preference for `hof env ci`: default, all, first
	// any are `hof env <key>`
	// sorta like hof/flow for hof/env

	// TODO, command for listening for webhooks and running jobs
	// still use argo events? (inside dagger?!)

	// we do want something short-term for
	// sequence / parallel / ephemeral for
	// tests / lint / ai check
	// abstracting to a env #pack / blueprint

	[string]: _
	// where do nested structs become nested lists, and do they come back again?
}

// we mainly need a way to  specify and/or track
// 1. pod like spec (resources, security primarily)
// 1. grouping and scheduling work, parallelism, sync
// 1. logical steps / error reporting
// 1. ci/cd graph, progress, logs, otel
// 1. final, consolidated report
