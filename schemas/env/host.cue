package env

// NOTE(2self): this probably needs a bookkeeping struct for { ...kind[<key>] }
// so we can track and reuse by kind and id, based on some "key" field, depending on the type

import (
	"github.com/hofstadter-io/hof/schemas"
)

// run a command on a host, only localhost for now
// WARNING, this does NOT go through dagger
// this is used in replacing ansible among other tools
// there is also an idea to have a flag that replaces the underlying runtime
//   such that [dagger,localhost,remote,kubernetes] becomes indistinguishable
// this is implemented with go.os/exec.Cmd, so mirrors it closely
#HostExec: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "hostExec"
	}

	$kind: "#hostExec"

	// the first arg is the Path, the rest are the args to it
	args: [string, ...string]

	// the working directory of the command
	// if not set, it is the current workdir hof is run from
	workdir?: string

	// key=value pairs
	envs: [...string]

	// filepath to redirect stdin to
	stdin?: string

	// filepath to redirect stdout to
	stdout?: string

	// filepath to redirect stderr to
	stderr?: string

	// expose all host env hof sees to the exec
	allEnv: bool | *false
}

// access an image in host container runtime
#HostImage: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "hostImage"
	}

	$kind: "#hostImage"

	// name of the image
	name: string
}

#HostFile: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "hostFile"
	}

	$kind: "#hostFile"

	// friendly name for file
	name?: string | *path

	// the path to load, relative or absolute
	path: string

	// a prefix to remove from the load path
	trimPrefix: string | *""

	// If true, the directory will always be reloaded from the host.
	noCache?: bool
}

#HostDir: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "hostDir"
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
	gitignore: bool | *true

	// a prefix to remove from the load path
	trimPrefix?: string

	// git-compatible patch to apply after selecting and trimming
	patch?:     string
	patchFile?: #FileLike
}

// Creates a service that forwards traffic to a specified address via the host.
// (proxy via the host?) (or is this how we expose: environ -> host)
#HostService: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "hostService"
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
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "hostTunnel"
	}

	$kind: "#hostTunnel"

	// friendly name for tunnel
	name?: string | service.hostname

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

// Accesses a Unix socket on the host.
#HostSocket: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "hostSocket"
	}

	$kind: "#hostSocket"

	// friendly name for socket
	name?: string | *path

	path: string
}

// #HostPatch

UnixSocket: Step & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "unixSocket"
	}

	$kind: "unixSocket"

	path: string
	source: #HostSocket

	owner?: string
	expand?: bool
}