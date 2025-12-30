package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

#GitRepo: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "gitRepo"
	}

	$kind: "#gitRepo"
	url:   string
	ref:   string | *"HEAD"

	name: string | *"\(url)@\(ref)"

	// opts
	keepGitDir:               bool | *true
	sshKnownHosts?:           string
	sshAuthSocket?:           #HostSocket
	httpAuthUsername?:        string
	httpAuthToken?:           #Secret
	httpAuthHeader?:          #Secret
	experimentalServiceHost?: #Service
}
