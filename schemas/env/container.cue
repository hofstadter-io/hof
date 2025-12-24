package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

// Definition for a container env
Container: {
	schemas.Hof
	#hof: env: {
		root: true
		kind: "container"
	}

	name: string

	// need some kind of from for host / git / oci
	from: string | Container

	steps: [...Step]
	labels: [string]: string

	labels: {
		"org.opencontainers.image.title":   string | *name
		"org.opencontainers.image.version": string | *"latest"
		"org.opencontainers.image.commit":  string | *"dirty"
	}
}

StepKinds: [
	// meta to force eval
	"sync",

	// running inside
	"exec",
	"name",
	"workdir",

	// filesys related
	"file",
	"getFile",
	"dir",
	"getDir",

	// envs & secrets
	"env",
	"envfile",
	"secret",

	// external resources
	"mount", // consolidated, we may want to split them?

	// runtime stuff
	"expose",
	"entrypoint",
	"args",
	"term", // default dagger term setup
]

Step: {
	$kind: or(StepKinds)
}

Sync: Step & {
	$kind: "sync"
}

Exec: Step & {
	$kind: "exec"
	args: [...string]

	// todo, think about how to hand stdio and redir to files,
	// ideally they can be on the CUE types, but this is when we get into...
	// the Fill CUE from Dagger results, continue eval'n CUE
	// OG w/ Dagger, we did this through flow (?), avoid that here if possible
	// just increasingly eval the value as much as we can?
	// figure out what still needs to happen in Dagger, then do that
	// this got more abstract than just Exec, applies to files/dir as well
}

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
	content?: string | File
}

// like dagger.File
#File: {
	$kind:   "#file"
	source?: string | Container | #Dir
	path:    string

	content: string
}

Dir: Step & {
	$kind:   "dir"
	path:    string
	source?: Container | #Dir
}

#Dir: {
	$kind:   "getDir"
	source?: string | Container | #Dir
	path:    string
	include: [...string]
	exclude: [...string]

	$out: #Dir | null
}

Env: Step & {
	$kind: "env"
	envs: [string]: string
}

Envfile: Step & {
	$kind: "envfile"
}

Secret: Step & {
	$kind: "secret"
}

Mount: Step & {
	$kind: "mount"
	// cache, dir, file, secret, temp, host, service (?)
}

Expose: Step & {
	$kind: "expose"
}

Entrypoint: Step & {
	$kind: "entrypoint"
	args: [...string]
}

Args: Step & {
	$kind: "args"
	args: [...string]
}

Term: Step & {
	$kind: "term"
	args: [...string]
}
