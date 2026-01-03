# `hof def ./schemas/env`

```cue
package env

// NOTE(2self): this probably needs a bookkeeping struct for { ...kind[<key>] }
// so we can track and reuse by kind and id, based on some "key" field, depending on the type

// hmmm, what might this be for...?

import "github.com/hofstadter-io/hof/schemas"

_cmdCommon: {
	// life-cycle, notifications, cleanup, ...
	hooks?: {
		onStart?:    _
		onProgress?: _
		onAbort?:    _
		onSuccess?:  _
		onFailure?:  _
	}

	// how to handle failures
	config?: {
		failFast: bool | *false
	}
}
#Cmd: {
	schemas.Hof
	_cmdCommon
	#hof: {
		env: {
			root: true
			kind: "cmd"
		}
	}
	$kind: "cmd"
	name:  string | *#hof.metadata.name

	// set tasks type and names
	tasks: {
		{
			[string]: #Task
		}
		{
			[k=string]: {
				name: k
				// steps: [...[...{name: "\(#hof.metadata.name).\(k)"}]]
			}
		}
	}
	...
}

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
#FileLike:  #File | #HostFile
#StepList: {
	$kind: !~"^#"
} | [...#StepList]
#HackList: [...] | {
	...
}

// Definition for a container env
#Container: {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "container"
		}
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
	envs: {
		[string]: string
	}

	// steps to build an image or environment
	// TODO, put some basic checking on this
	steps: [...]

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
EnvVar: Step & {
	{
		$kind:    "envVar"
		[string]: string
	}
}
EnvFile: Step & {
	$kind: "envFile"
	file:  #File | #HostFile
}

// pass all os.Env vars to the container or service
EnvAll: Step & {
	$kind: "envAll"
}

// sets a secret in the system
#Secret: Step & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "secret"
		}
	}
	$kind: "#secret"
	name:  string

	// plaintext, uri, or file
	// actual, import env/rrr:env to enforce, performance penalty included
	// source: string | #FileLike
	source: _
}
SecretVar: Step & {
	$kind: "secretVar"

	// the secret VAR_NAME
	name: string

	// the secret value
	secret: #Secret
}
Exec: Step & {
	$kind: "exec"
	args: [...string]
	useEntrypoint?:                 bool
	stdin?:                         string
	redirectStdin?:                 string
	redirectStdout?:                string
	redirectStderr?:                string
	expect?:                        *"SUCCESS" | "FAILURE" | "ANY"
	experimentalPrivilegedNesting?: bool
	insecureRootCapabilities?:      bool
	expand?:                        bool
	noInit?:                        bool
}
Script: Exec & {
	script: string
	args: ["sh", "-c", script]
}
Sh: Exec & {
	script: string
	_script: """
		set -euo pipefail

		"""
	args: ["sh", "-c", _script + script]
}
Bash: Exec & {
	script: string
	_script: """
		set -euo pipefail

		"""
	args: ["bash", "-c", _script + script]
}
Zsh: Exec & {
	script: string
	_script: """
		set -euo pipefail

		"""
	args: ["zsh", "-c", _script + script]
}
Sync: Step & {
	$kind: "sync"
}
User: Step & {
	$kind: "user"
	name:  string
}
Workdir: Step & {
	$kind: "workdir"
	path:  string
}
Entrypoint: Step & {
	$kind: "entrypoint"
	args: [...string]
	keepDefaultArgs?: bool
}
DefaultArgs: Step & {
	$kind: "defaultArgs"
	args: [...string]
}

// sets the default terminal
DefaultTerm: Step & {
	$kind: "defaultTerm"
	args: [...string]
	experimentalPrivilegedNesting?: bool
	insecureRootCapabilities?:      bool
}
#ExportDir: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "exportDir"
		}
	}
	$kind: "#exportDir"
	name:  string

	// where to place
	path: string

	// pieces that make up the bundled dir
	// sources:  [...#FileLike|#DirLike]
	sources: [...]

	// (1) filters
	include: [...string]
	exclude: [...string]
	gitignore: bool | *true

	// (2) path to select from the bundled dir
	bundlePath?: string

	// (3) git-compatible patch to apply after bundling and selecting
	patch?:     string
	patchFile?: #FileLike

	// If true, then the host directory will be wiped clean before exporting so that it exactly matches the directory being exported; this means it will delete any files on the host that aren't in the exported dir. If false (the default), the contents of the directory will be merged with any existing contents of the host directory, leaving any existing files on the host that aren't in the exported directory alone.
	wipe: bool | *false
}
#ExportFile: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "exportFile"
		}
	}
	$kind: "#exportFile"
	name:  string
	path:  string
	file:  #File

	// If allowParentDirPath is true, the path argument can be a directory path, in which case the file will be created in that directory.
	allowParentDirPath?: bool
}
#ExportImageFile: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "exportImageFile"
		}
	}
	$kind: "#exportImageFile"
	name:  string
	path:  string
	tags: [...string]
	image: #Container
}
#ExportImage: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "exportImage"
		}
	}
	$kind: "#exportImage"
	name:  string
	reg?:  string
	tags: [...string]
	image: #Container
}
#PublishImage: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "publishImage"
		}
	}
	$kind: "#publishImage"
	name:  string
	reg:   string
	tags: [...string]
	image: #Container
}

// export the CUE representation
#ExportCuefig: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "exportCuefig"
		}
	}
	$kind:   "#exportCuefig"
	name:    string
	path:    string
	format?: string
	data:    _
}

// like dagger.File
#File: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "file"
		}
	}
	$kind: "#file"
	name:  string | *path
	path!: string

	// actual, import env/rrr:env to enforce, performance penalty included
	// source: #Dir | #Container | #HostDir | #HostImage | #GitRepo
	source!: _
}

// this is creating a directory ref that we can do things with
#Dir: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "dir"
		}
	}
	$kind: "#dir"
	name:  string | *path

	// where to place
	path: string | *"."

	// pieces that make up the bundled dir
	// sources:  [...#FileLike|#DirLike]
	sources: [...]

	// (1) filters
	include: [...string]
	exclude: [...string]
	gitignore: bool | *true

	// (2) path to select from the bundled dir
	bundlePath: string | *"/"

	// (3) git-compatible patch to apply after bundling and selecting
	patch?:     string
	patchFile?: #FileLike
}

// like dagger.WithFile
File: Step & {
	$kind: "file"
	path!: string

	// actual, import env/rrr:env to enforce, performance penalty included
	// content: string | #File | #HostFile // HMMM(A): should this just be file, or be container/image too?
	content!:     _
	permissions?: int
	owner?:       string
	expand?:      bool
}
#GitRepo: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "gitRepo"
		}
	}
	$kind: "#gitRepo"
	url:   string
	ref?:  string
	name:  string | *"\(url)@\(*ref | "HEAD")"

	// opts
	keepGitDir:               bool | *true
	sshKnownHosts?:           string
	sshAuthSocket?:           #HostSocket
	httpAuthUsername?:        string
	httpAuthToken?:           #Secret
	httpAuthHeader?:          #Secret
	experimentalServiceHost?: #Service
}

// access an image in host container runtime
#HostImage: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "hostImage"
		}
	}
	$kind: "#hostImage"

	// name of the image
	name: string
}
#HostFile: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "hostFile"
		}
	}
	$kind: "#hostFile"

	// friendly name for file
	name?: string | *path

	// the path to load, relative or absolute
	path: string

	// If true, the directory will always be reloaded from the host.
	noCache?: bool
}
#HostDir: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "hostDir"
		}
	}
	$kind: "#hostDir"

	// friendly name for dir
	name: string | *path

	// the path to load, relative or absolute
	path!: string

	// Exclude artifacts that match the given pattern (e.g., ["node_modules/", ".git*"]).
	exclude?: [...string]

	// Include only artifacts that match the given pattern (e.g., ["app/", "package.*"]).
	include?: [...string]

	// If true, the directory will always be reloaded from the host.
	noCache?: bool

	// Apply .gitignore filter rules inside the directory
	gitignore?: bool | *true
}

// Creates a service that forwards traffic to a specified address via the host.
// (proxy via the host?) (or is this how we expose: environ -> host)
#HostService: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "hostService"
		}
	}
	$kind: "#hostService"

	// friendly name for service
	name?: string | *host

	// Upstream host to forward traffic to.
	host: string | *"localhost"

	// Configure explicit port forwarding rules for the service.
	// If a port's frontend is unspecified or 0, a random port will be chosen by the host.
	// If no ports are given, all of the service's ports are forwarded. If native is true, each port maps to the same port on the host. If native is false, each port maps to a random port chosen by the host.
	// If ports are given and native is true, the ports are additive.
	ports?: [...#PortForward]
}

// Creates a tunnel that forwards traffic from the host to a service.
// (host -> environ)
#HostTunnel: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "hostTunnel"
		}
	}
	$kind: "#hostTunnel"

	// friendly name for tunnel
	name?:   string | service_9.hostname
	service: #Service

	// Map each service port to the same port on the host, as if the service were running natively.
	// Note: enabling may result in port conflicts.
	native?: bool

	// Configure explicit port forwarding rules for the tunnel.
	// If a port's frontend is unspecified or 0, a random port will be chosen by the host.
	// If no ports are given, all of the service's ports are forwarded. If native is true, each port maps to the same port on the host. If native is false, each port maps to a random port chosen by the host.
	// If ports are given and native is true, the ports are additive.
	ports?: [...#PortForward]
}

// something you can launch or deploy, this amounts to AsService in Dagger
#Service: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "service"
		}
	}
	$kind: "#service"

	// convenient name
	name!: string

	// container to turn into a service
	// actual, import env/rrr:env to enforce, performance penalty included
	// source: #Container | #HostImage
	source!: _

	// ports to expose on the container
	ports?: [...#PortForward]

	// configures a hostname within the session at which the server which it can be reached
	// used when exposing to the host
	hostname: string | *name

	// if empty, the container's default will be used
	args?: [...string]

	// todo
	envVars?: {
		[string]: string
	}
	envFiles?: {
		[string]: #File
	}
	shhVars?: {
		[string]: string
	}
	shhFiles?: {
		[string]: #File
	}

	// if the container has an entrypoint, prepend it to the args
	useEntrypoint?: bool

	// Provides Dagger access to the executed command.
	experimentalPrivilegedNesting?: bool

	// Execute the command with all root capabilities. This is similar to running a command with "sudo" or executing "docker run" with the "--privileged" flag. Containerization does not provide any security guarantees when using this option. It should only be used when absolutely necessary and only with trusted commands.
	insecureRootCapabilities?: bool

	// Replace "${VAR}" or "$VAR" in the args according to the current environment variables defined in the container (e.g. "/$VAR/foo").
	expand?: bool

	// If set, skip the automatic init process injected into containers by default.
	//
	// This should only be used if the user requires that their exec process be the pid 1 process in the container. Otherwise it may result in unexpected behavior.
	noInit?: bool
}
#PortForward: {
	// friendly name for the port
	name?: string
	port:  int

	// Destination port for traffic.
	backend: port

	// Port to expose to clients. If unspecified, a default will be chosen.
	frontend?: int

	// Transport layer protocol to use for traffic.
	protocol: *"TCP" | "UDP"
}
Expose: Step & {
	$kind:                         "expose"
	name?:                         string
	port:                          int
	protocol:                      *"TCP" | "UDP"
	experimentalSkipHealthchecks?: bool
}
#Space: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "space"
		}
	}
	$kind:     "space"
	name:      string
	terminal?: #Container
	services?: {
		[string]: #Service
	}
	external?: {
		[string]: _
	}

	// tbd...
	cmds?: {
		[string]: _
	}
	flags?: {
		[string]: _
	}
	config?: {
		[string]: _
	}
	stacks?: {
		[string]: _
	}
	exports?: {
		[string]: _
	}
}
#Cache: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "cache"
		}
	}
	$kind: "#cache"
	name:  string
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
#Task: {
	schemas.Hof
	_cmdCommon
	#hof: {
		env: {
			root: true
			kind: "task"
		}
	}
	$kind: "task"
	name:  string

	// ideally, this is more dag/flow like
	// two-level list, top-sequential | nest-parallel
	// TODO, put some basic checking on this
	steps: [...[...]]
	...
}
#DockerBuild: {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "dockerBuild"
		}
	}
	$kind:       "#dockerBuild"
	name?:       string
	source:      #Dir | #HostDir
	dockerfile?: string
	platform?:   string
	buildArgs: {
		[string]: string
	}
	target?: string
	secrets?: [...#Secret]
	noInit?: bool
}

// treat secret content is an env file
// exposing each line as secret vars
SecretFile: Step & {
	$kind:  "secretFile"
	source: #Secret
}

// starts an interactive terminal
Terminal: Step & {
	$kind: "terminal"
	args: [...string]
	experimentalPrivilegedNesting?: bool
	insecureRootCapabilities?:      bool
}

// export the Dagger representation
#ExportDagger: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "exportDagger"
		}
	}
	$kind:   "#exportDagger"
	name:    string
	path:    string
	format?: string
	data:    _
}

// this is including a directory in a container
Dir: Step & {
	$kind: "dir"
	// args
	path: string | *"."

	// actual, import env/rrr:env to enforce, performance penalty included
	// source: #Container | #Dir | #GitRepo | #HostDir | #HostImage // HMMM(B): or maybe this should just be dir kinds, make the user do an extra step? (nah, wouldn't have to with the SDK directly)
	source!: _

	// opts
	include?: [...string]
	exclude?: [...string]
	gitignore?: bool | *true
	owner?:     string
	expand?:    bool
}

// Accesses a Unix socket on the host.
#HostSocket: Ref & {
	schemas.Hof
	#hof: {
		env: {
			root: true
			kind: "hostSocket"
		}
	}
	$kind: "#hostSocket"

	// friendly name for socket
	name?: string | *path
	path:  string
}
BindService: Step & {
	$kind: "bindService"

	// confitures an alias for the service when binding to this container
	alias:   string | *self.service.name
	service: #Service
}
Mount: Step & {
	$kind: "mount"
	path:  string

	// cache, dir, file, secret, temp, host, service (?)
	// source: #Cache | #File | #HostFile | #Dir | #HostDir
	source: _
}

let service_9 = service
```