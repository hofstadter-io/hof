---
title: "Modules"
description: "Dependency management for CUE and code generators"
brief: "Dependency management for CUE and code generators"

weight: 20
---

### Hof & CUE Modules

Every __hof generator__ is also a __CUE module__,
and in fact, many of hof's other features can
be used from the module system.

"`hof mod`" is the subcommand based on Go modules
for working with CUE modules and managing dependencies.
The logic and format is the same, with much code shared between the implementations.
Currently, while CUE is module aware and supports imports,
it does not yet have dependency management, but it will work the same as well.
This page has a brief overview. The commands will
be repeated when you need to run them on other pages.

`hof mod` subcommands:

- `hof mod init cue github.com/user/repo` initializes a new module
- `hof mod vendor cue` fetches dependencies into `cue.mod/pkg/...`

The name of a module must be the same the git repository.
`hof` talks directly to git repositiories and many of
`hof`'s commands will accept repositories as input too.

The files and directories that make up a module:

- `cue.mods` is where dependencies and versions are set, you write this file.
- `cue.sums` contains the checksums for all dependencies and is managed by `hof mod`
- `cue.mod/module.cue` denotes a CUE module and has a sinlge line
- `cue.mod/pkg/...` is where the code for dependencies is located after fetching

{{<codeInner title="cue.mod/module.cue">}}
// indicates a CUE module, only one line
module: "github.com/user/repo"
{{</codeInner>}}

{{<codeInner title="cue.mods">}}
// indicates a Hof module, dependencies go here
module github.com/user/repo

cue {{<cue-version>}}

require (
	github.com/hofstadter-io/hof {{<hof-version>}}
)
{{</codeInner>}}


### Replace for local development

You can use `hof mod` for generally managing CUE or Hof modules and dependencies.
In addition to setting, fetching, and validating dependencies,
you can use `hof mod` to setup local development when working
with multiple CUE or Hof modules.

In your `cue.mods` file, use the `replace` directive just like you would in Go.

{{<codeInner title="cue.mods with replace">}}
module github.com/user/repo

cue {{<cue-version>}}

require (
	github.com/username/repo v0.1.2
)

// local replace with relative path
replace github.com/username/repo => ../repo
{{</codeInner>}}

{{<codeInner title="> terminal">}}
$ hof mod vendor cue
{{</codeInner>}}

This will symlink any local replaces
to point from the CUE vendor directory
to the replacement directory.
You can then develop CUE or hof code
without having to copy or vendor upstream dependencies.

