package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

// Step in #container: [...Step]
Step: {
	$kind: string
}

// #Things...
Ref: {
	$kind: string
	id?:   string
}

#ImageLike: #Container | #HostImage | #DockerBuild
#DirLike:   #Dir | #HostDir | #GitRepo

#FileLike: #File | #HostFile

#StepList: {$kind: string & !~"^#"} | [...#StepList]

#HackList: [...] | {...}

// todo, registry auth

// Definition for a container env
#Container: {
	// always configure a veg node around containers
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "container"
	}

	// duplicative, but needed for the decoding switch statement simplicity (it uses $kind, while `veg node` uses #hof.kind)
	$kind: "#container"

	// the name of the container or environment
	name?: string

	// actual, import env/rrr:env to enforce, performance penalty included
	// from: string | #Container | #HostImage | #DockerBuild
	// from: string | {...}  !!! PANIC !!!
	from!: _

	// you can do this in steps, but it might be nice to have it
	// 1. extracted / separate for easy usage in k8s (i.e.)
	// 2. auto add them at the beginning before any steps
	// this seems a reasonable DX
	envs: [string]: string

	// steps to build an image or environment
	// TODO, put some basic checking on this
	steps: [...] // OK

	// steps: [...#HackList] // PANIC
	// steps: [...([...] | {...})] // PANIC
	// steps: [...] | {...} // OK (but not right)
	// steps: [...{...}]  // OK (but not right)
	// steps: [...{...}|[...]]  // PANIC
}

DefaultLabels: {
	#name?:                             string | *"ephemeral"
	"org.opencontainers.image.title":   string | *#name
	"org.opencontainers.image.version": string | *"latest"
	"org.opencontainers.image.commit":  string | *"dirty"
}

#DockerBuild: {
	schemas.Hof
	#hof: env: {
		root: true
		kind: "dockerBuild"
	}

	$kind: "#dockerBuild"
	name?: string

	source: #Dir | #HostDir

	dockerfile?: string
	platform?:   string
	buildArgs: [string]: string

	target?: string
	secrets?: [...#Secret]
	noInit?: bool
}

// we probably need to move this into the Go
// so we can copy over a bunch of the meta/env/cmd/entry
#Flatten: #Container & {
	#orig: _
	from: "scratch"
	steps: [
		Dir & { path: "/", source: #Dir & { path: "/", sources: [#orig]} }
	]
}