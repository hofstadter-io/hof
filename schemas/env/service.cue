@experiment(aliasv2)
package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

// something you can launch or deploy, this amounts to AsService in Dagger
#Service: Ref & {
	schemas.Hof
	#hof: env: {
		root: true
		kind: "service"
	}

	$kind: "#service"

	// convenient name
	name?: string

	// container to turn into a service
	// actual, import env/rrr:env to enforce, performance penalty included
	// source: #Container | #HostImage
	source!: _

	// ports to expose on the container
	ports?: [...#PortForward]

	// configures a hostname within the session at which the server which it can be reached
	// used when exposing to the host
	hostname?: string

	// if empty, the container's default will be used
	args?: [...string]

	// todo
	envVars?:  [string]: string
	envFiles?: [string]: #File
	shhVars?:  [string]: string
	shhFiles?: [string]: #File


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

	port: int

	// Destination port for traffic.
	backend: port
 
	// Port to expose to clients. If unspecified, a default will be chosen.
	frontend?: int

	// Transport layer protocol to use for traffic.
	protocol: *"TCP" | "UDP"
}

Expose: Step & {
	$kind: "expose"

	name?:    string
	port:     int
	protocol: *"TCP" | "UDP"

	experimentalSkipHealthchecks?: bool
}

BindService: Step & {
	$kind: "bindService"

	// confitures an alias for the service when binding to this container
	alias?:   string | *self.service.hostname
	service: #Service
}
