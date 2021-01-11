# hof - the high code framework

The `hof` tool tries to remove redundent development activities
by using high level designs, code generation, and diff3
while letting you write custom code directly in the output.
( low-code for developers )

- Users write Single Source of Truth (SSoT) design for data models and the application generators
- `hof` reads the SSoT, processes it through the code generators, and outputs directories and files
- Users can write custom code in the output, change their designs, and regenerate code in any order
- `hof` can be customized and extended by only editing text files and not `hof` source code.
- Use your own tools, technologies, and practices, `hof` does not make any choices for you
- `hof` is powered by Cue (https://cuelang.org & https://cuetorials.com)

## Install

You will have to download `hof` the first time.
After that `hof` will prompt you to update and
install new releases as they become available.

```text
# Install (Linux, Mac, Windows)
curl -LO https://github.com/hofstadter-io/hof/releases/download/v0.5.15/hof_0.5.15_$(uname)_$(uname -m)
mv hof_0.5.15_$(uname)_$(uname -m) /usr/local/bin/hof

# Shell Completions (bash, zsh, fish, power-shell)
echo ". <(hof completion bash)" >> $HOME/.profile
source $HOME/.profile

# Show the help text
hof --help
```

You can always find the latest version from the
[releases page](https://github.com/hofstadter-io/hof/releases)
or use `hof` to install a specific version of itself with `hof update --version vX.Y.Z`.


## Documentation

Please see __https://docs.hofstadter.io__ to learn more.

Join us on Slack! https://hofstadter-io.slack.com [(invite link)](https://join.slack.com/t/hofstadter-io/shared_invite/zt-e5f90lmq-u695eJur0zE~AG~njNlT1A)


## Example

There are currently hof modules for:

- [hofmod-cli](https://github.com/hofstadter-io/hofmod-cli) - CLI infrastructure based on the Golang Cobra library.
- [hofmod-server](https://github.com/hofstadter-io/hofmod-server) - API server based on the Golang Echo library.

You can see them used in:

-  `hof` uses `hofmod-cli`
- [saas](https://github.com/hofstadter-io/_saas) uses `hofmod-server` 

The following is a single file example:


```
package gen

import (
	// import hof's schemas for our generator
	"github.com/hofstadter-io/hof/schema"
)

// A schema for our generator's input
#Input: {
	name: string
	todos: [...{
		name: string
		effort: int
		complete: bool
	}]
}
// create a generator
#Gen: schema.#HofGenerator & {
	// We often have some input values for the user to provide.
	// Use a Cue definition to enforce a schema
	Input: #Input

	// Required filed for generator definitions, details can be found in the hof docs
	PackageName: "dummy"

	// Required field for a generator to work, the list of files to generate
	Out: [...schema.#HofGeneratorFile] & [
		todo,
		done,
		debug,
	]

	// In is supplied as the root data object to every template
	// pass user inputs to the tempaltes here, possibly modified, enhanced, or transformed
	In: {
		INPUT: Input
		Completed: _C
		Incomplete: _I
	}

	// calculate some internal data from the input
	_C: [ for t in Input.todos if t.complete == true { t } ]
	_I: [ for t in Input.todos if t.complete == false { t } ]

	// the template files
	todo: {
		Template: """
		Hello {{ .INPUT.name }}.

		The items still on your todo list:

		{{ range $T := .Incomplete -}}
		{{ printf "%-4s%v" $T.name $T.effort }}
		{{ end }}
		"""
		// The output filename, using string interpolation
		Filepath: "\(Input.name)-todo.txt"
	}
	done: {
		Template: """
		Here's what you have finished {{ .INPUT.name }}. Good job!

		{{ range $T := .Completed -}}
		{{ $T.name }}
		{{ end }}
		"""
		Filepath: "\(Input.name)-done.txt"
	}

	// useful helper
	debug: {
		Template: """
		{{ yaml . }}
		"""
		Filepath: "debug.yaml"
	}
}

// Add the @gen(<name>,<name>,...) to denote usage of a generator
Gen: _ @gen(todos)
// Construct the generator
Gen: #Gen & {
	Input: {
		// from first.cue
		name: gen.data.name
		todos: gen.data.tasks
	}
}
```
