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
}

// temp space config for ephemeral volumes not persisted between exec calls
#Temp: {
	$kind: "#temp"

	// where to attach it
	path: string

	// size in bytes
	size?: int

	// expand vars in path like $HOME/.cache
	expand?: bool
}
