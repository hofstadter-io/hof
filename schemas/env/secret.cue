package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

// sets a secret in the system
#Secret: Step & {
	schemas.Hof
	#hof: env: {
		root: true
		kind: "secret"
	}

	$kind: "#secret"
	name: string

	// plaintext, uri, or file
  // actual, import env/rrr:env to enforce, performance penalty included
	// source: string | #FileLike
  source: _
}

Secret: Step & {
	$kind: "secret"

  // the secret VAR_NAME
  name: string

  // the secret value
  secret: #Secret
}

// treat secret content is an env file
// exposing each line as secret vars
Secretvars: Step & {
	$kind: "secretvars"

	source: #Secret
}
