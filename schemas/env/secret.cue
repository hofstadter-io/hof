package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

// Secret: Step & {
// 	$kind:  "secret"
// 	var?:   string
// 	secret: #Secret
// }

#Secret: Step & {
	schemas.Hof
	#hof: env: {
		root: true
		kind: "secret"
	}

	$kind: "#secret"

	name: string

	// plaintext, uri, or file
	source: string | #File | #HostFile

	owner?:  string
	expand?: bool
	mode?:   int
}
