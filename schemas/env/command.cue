@experiment(aliasv2)

package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

#Cmd: {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "cmd"
	}

	$kind: "cmd"

	name: string

	// ideally, this is more dag/flow like
	// two-level list, top-sequential | nest-parallel
	tasks: [string]~(k,_): #Task & {name: k}

	// how to handle failures
	config?: {
		failFast: bool | *false
	}

	// life-cycle, notifications, cleanup, ...
	hooks?: {
		onStart:    _
		onProgress: _
		onAbort:    _
		onSuccess:  _
		onFailure:  _
	}

	...
}

#Task: {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "task"
	}

	$kind: "task"
	name:  string

	parallel: bool | *false
	tasks: [...#Task]
}
