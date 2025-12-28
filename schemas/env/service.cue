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

  // convenience, and the default for hostname/alias depending on where it is used
  name: string | *hostname

  // configures a hostname within the session at which the server which it can be reached
  hostname: string | *name

  // ports to expose on the container
  ports?: [...#Port]

  // container to turn into a service
  source: #Container | #HostImage

  // if empty, the container's default will be used
  args?: [...string]

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

#Port: {
  name?: string
  proto: *"tcp" | "udp" // align this with k8s too

  // the port
  port: int
  hostPort: int | *port

	experimentalSkipHealthchecks?: bool
}