/* package runtimes

Runtimes are Cue values with a standard
form and configuration that Hofstadter
understands, implements, and interleaves
with your Cue code.

Runtimes enable you to write scripts and imperative code
in many common languages and existing runtimes:

```
foo: "bar"

// dev todo, think about runtime specializations
// due to bash not having a nice way to nest fields / structs and access?
// and others do? Or can we do this with runtime CRDs that have a name and an entrypoint / args
// we plan to exec out to mode, can probably text interpolate the code, getting it back in is harder
// and probably up to what the user is doing

// NOTE, we may be (can we make it an option) writing code to create cue files
// which are exec'd in order, but are still valid cue files on their own (and when eval'd by cue itself)?

toUpper: hof.#Runtime @runtime(bash)
toUpper:{
	// top-level fields are defined
	Args: {
		name: string
	}
	Code: """
	echo ${name^^}
	"""
	Result: {
		NAME: string @stdout()
		Error: string ^stderr()
	}
}

FOO:
```

*/

package runtimes
