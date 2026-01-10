@experiment(aliasv2)

package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

_cmdCommon: {
	// life-cycle, notifications, cleanup, ...
	hooks?: {
		onStart?:    _
		onProgress?: _
		onAbort?:    _
		onSuccess?:  _
		onFailure?:  _
	}

	// how to handle failures
	config?: {
		failFast: bool | *false
	}
}

#Cmd: {
	@env()
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "cmd"
	}

	$kind: "cmd"
	name:  string | *#hof.metadata.name

	// set tasks type and names
	tasks: [string]: #Task
	tasks: [string]~(k,_): {
		@env()

		name: k
		// steps: [...[...{name: "\(#hof.metadata.name).\(k)"}]]
	}

	_cmdCommon

	...
}

#Task: {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "task"
	}

	$kind: "task"
	name:  string | *#hof.metadata.name

	// ideally, this is more dag/flow like
	steps: [...]

	parallel: int | *0

	_cmdCommon

	...
}
