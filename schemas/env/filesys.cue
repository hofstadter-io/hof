package env

// like dagger.File
#File: Ref & {
	$kind:   "#file"
	path:    string
	source?: #Dir | #Container | #HostDir | #HostImage
}

// this is creating a directory ref that we can do things with
#Dir: Ref & {
	$kind:   "#dir"
	source?: #Dir | #Container | #HostDir | #HostImage
	path:    string
	include: [...string]
	exclude: [...string]
	gitignore: bool | *true
}

#Git: Ref & {
	$kind: "#git"
	url:   string

	// opts
	keepGitDir?:              bool | *true
	sshKnownHosts?:           string
	sshAuthSocket?:           #HostSocket
	httpAuthUsername?:        string
	httpAuthToken?:           #Secret
	httpAuthHeader?:          #Secret
	experimentalServiceHost?: #Service
}
