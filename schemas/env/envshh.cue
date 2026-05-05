package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

EnvVars: Step & {
	$kind:    "envVars"
	[string]: string
}

EnvFile: Step & {
	$kind: "envFile"

	file: #File | #HostFile
}

// pass all os.Env vars to the container or service
EnvAll: Step & {
	$kind: "envAll"
}

// sets a secret in the system
#Secret: Step & {
	schemas.Hof
	#hof: env: {
		root: true
		kind: "secret"
	}

	$kind: "#secret"
	name:  string

	// plaintext, uri, or file
	// actual, import env/rrr:env to enforce, performance penalty included
	// source: string | #FileLike
	source: _
}

SecretVars: Step & {
	$kind: "secretVars"

	// the secret value
	[!~#"\$kind"#]: #Secret
}

// treat secret content is an env file
// exposing each line as secret vars
SecretFile: Step & {
	$kind: "secretFile"

	file: #File | #HostFile | #Secret
}
