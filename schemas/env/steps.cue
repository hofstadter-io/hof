package env

Step: {
	$kind: or(StepKinds)
}

Ref: Step & {
	id?: string
}

Sync: Step & {
	$kind: "sync"
}

Exec: Step & {
	$kind: "exec"
	args: [...string]

	useEntrypoint?:  bool
	stdin?:          string
	redirectStdin?:  string
	redirectStdout?: string
	redirectStderr?: string
  expect?: *"SUCCESS" | "FAILURE" | "ANY"

	experimentalPrivilegedNesting?: bool
	insecureRootCapabilities?:      bool
	expand?:                        bool
	noInit?:                        bool
}

// todo, think about how to hand stdio and redir to files,
// ideally they can be on the CUE types, but this is when we get into...
// the Fill CUE from Dagger results, continue eval'n CUE
// OG w/ Dagger, we did this through flow (?), avoid that here if possible
// just increasingly eval the value as much as we can?
// figure out what still needs to happen in Dagger, then do that
// this got more abstract than just Exec, applies to files/dir as well

User: Step & {
	$kind: "user"
	name:  string
}

Workdir: Step & {
	$kind: "workdir"
	path:  string
}

// like dagger.WithFile
File: Step & {
	$kind: "file"

	path: string

	// what if from another container? helper to make one of these from a Container?
	content?: string | #File | #HostFile // HMMM(A): should this just be file, or be container/image too?

	permissions?: int
	owner?:       string
	expand?:      bool
}

// this is including a directory in a container
Dir: Step & {
	$kind: "dir"
	// args
	path:    string
	source?: #Container | #Dir | #HostDir | #HostImage // HMMM(B): or maybe this should just be dir kinds, make the user do an extra step? (nah, wouldn't have to with the SDK directly)
	// opts
	include?: [...string]
	exclude?: [...string]
	gitignore?: bool | *true
	owner?:     string
	expand?:    bool
}

Mount: Step & {
	$kind: "mount"

	path: string

	// cache, dir, file, secret, temp, host, service (?)
	source?: #Cache | #Dir | #File | #HostDir | #HostFile
}

Env: Step & {
	$kind:   "env"
  // bit of a hack for convenience in a couple places
	$expand: "1" | "t" | "T" | "TRUE" | "true" | "True" | "0" | "f" | "F" | "FALSE" | "false" | *"False"
  [string]: string
}

Envfile: Step & {
	$kind: "envfile"

	file: #File
}

Secret: Step & {
	$kind:  "secret"
	var?:   string
	secret: #Secret
}

Expose: Step & {
	$kind: "expose"

	name:     string
	port:     int
	protocol: *"tcp" | "udp"

	experimentalSkipHealthchecks?: bool
}

BindService: Step & {
	$kind: "bindService"

	// confitures an alias for the service when binding to this container
	alias:   string | *service.hostname
	service: #Service
}

Entrypoint: Step & {
	$kind: "entrypoint"
	args: [...string]

	keepDefaultArgs?: bool
}

Args: Step & {
	$kind: "args"
	args: [...string]
}

Term: Step & {
	$kind: "term"
	args: [...string]

	experimentalPrivilegedNesting?: bool
	insecureRootCapabilities?:      bool
}
