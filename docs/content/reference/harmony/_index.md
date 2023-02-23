---
title: "harmony"
weight: 100
---

{{<lead>}}
`harmony` is a framework for
testing a suite of repositories
across versions of dependencies.
Discover errors and regressions
in downstream projects before
releasing new code.
{{</lead>}}

`harmony` makes it easy for both developers and users
to setup and register projects for usage in upstream testing.
Easily test against any commit or working code.
Easily register projects with a small amount of CUE or config.

Built on [Dagger](https://dagger.io), `harmony` provides

- simple setup for developers
- simple registration for user projects
- support for any languages or tools 
- consistent environment for testing
- simplified version injection
- easily run specific or all cases

_Note, while `harmony` uses Dagger, it is not required for downstream projects._

#### Examples

- [harmony-cue](https://github.com/hofstadter-io/harmony-cue) for testing CUE base projects

## Harmony Setup

`harmony` is just a Dagger Plan, so to setup your project
you only need to add the following file to your project.

```cue
package main

import (
  "dagger.io/dagger"
  "universe.dagger.io/docker"

  "github.com/hofstadter-io/harmony"

  // import your projects registry
  "github.com/username/project/registry"
)

// A dagger plan is used as the driver for testing
dagger.#Plan

// add actions from Harmony
actions: harmony.Harmony

// project specific actions & configuration
actions: {

  // global version config for this harmony
  versions: {
    go: "1.18"
  }

  // the registry of downstream projects
  // typically we put this in a subdir and import it
  "registry": registry.Registry

  // the image test cases are run in
  // typically parametrized so we can change dependencies or versions 
  // (you can also build any image you want here)
  runner: docker.#Pull & {
    source: "index.docker.io/golang:\(versions.go)-alpine"
  }

  // where downstream project code is checked out
  workdir: "/work" 
}
```

Run registration cases or a single case with dagger: `dagger do <reg> [case]`.
Any cases found will be run in parallel.

Use `./run.py` to run the full suite of registrations and cases sequentially,
or as a convenient way to set dependency versions.

#### Registry

You will typically want to provide a subdirectory
for registered projects. You can also provide
short codes to simplify user project registration further.

Here we add a `_dagger` short code by
including a `schema.cue` in our `registry/` directory.

```cue
package registry

import (
  "strings"

  "universe.dagger.io/docker"
  "github.com/hofstadter-io/harmony"
)

// customized Registration schema built on harmony's
Registration: harmony.Registration & {
  // add our short codes 
  cases: [string]: docker.#Run & {
    _dagger?: string
    if _dagger != _|_ {
      command: {
        name: "bash"
        args: ["-c", _script]
        _script: """
        dagger version
        dagger project update
        dagger do \(_dagger)
        """ 
      }
    }
  }
}
```

Registrations then use with `cases: foo: { _dagger: "foo bar", workdir: "/work" }`


## Registration Setup

To add a new registration for user projects,
add a CUE file to the upstream project.

```cue
package registry

// Note the 'Registry: <name>: ...` needs to be unique
Registry: hof: Registration & {
  // 
  remote: "github.com/hofstadter-io/hof"
  ref: "main"

  cases: {
    // these are docker.#Run from Dagger universe
    cli: { 
      workdir: "/work"
      command: {
        name: "go"
        args: ["build", "./cmd/hof"]
    }
    flow: { 
      workdir: "/work"
      command: {
        name: "go"
        args: ["test", "./flow"]
    }
  }
}
```

## Base image using dirty code

As a developer of a project,
you may wish to run harmony
before committing your code.
To do this, we need to

- mount the local directory into dagger
- build a runner image from the local code
- mount the local directory into the runner (possibly)

Refer to [harmony-cue](https://github.com/hofstadter-io/harmony-cue) for an example.
Look at how "local" is used:

- harmony.cue
- testers/image.cue
- run.py

The essence is to use CUE guards (if statements)

---

`harmony` was inspired by the idea
of creating a [cue-unity](https://github.com/cue-unity/unity)
powered by [Dagger](https://dagger.io).
The goal of `cue-unity` is to collect community projects
and run them against new CUE language changes.
`cue-unity` was itself inspired by some work
Rob Pike did on Go, for the same purpose
of testing downstream projects using Go's stdlib.

`harmony` is more generalized and
the CUE specific version is [harmony-cue](https://github.com/hofstadter-io/harmony-cue).

