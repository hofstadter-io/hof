package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

#Cache: Ref & {
	schemas.Hof
	#hof: env: {
		root: true
		kind: "cache"
	}

	$kind: "#cache"
	name:  string

	// dir-like to prepopulate cache with
	source?: _

	// watch the source #HostDir for changes
	// only works with #HostDir and Mount
	watch?: bool
}

// temp space config for ephemeral volumes not persisted between exec calls
Temp: {
	$kind: "temp"

	// where to attach it
	path: string

	// size in bytes
	size?: int

	// expand vars in path like $HOME/.cache
	expand?: bool
}

Mount: Step & {
	$kind: "mount"

	path: string

	// cache, dir, file, secret, temp, host, service (?)
	// source: #Cache | #File | #HostFile | #Dir | #HostDir
	source: _
	expand: bool | *true

	// cache, file, dir, secret
	owner?: string
	// secret
	mode?: int
}
