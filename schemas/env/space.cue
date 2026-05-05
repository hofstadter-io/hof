package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

#Space: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "space"
	}

	$kind: "space"
	name:  string

	//
	// collection of all the other things
	//

	terminal?: #Container
	services?: [string]: #Service
	external?: [string]: _ // stuff from host or git, mainly filesystem like stuff

	// tbd...
	cmds?: [string]:    _
	flags?: [string]:   _
	config?: [string]:  _
	stacks?: [string]:  _
	exports?: [string]: _
}
