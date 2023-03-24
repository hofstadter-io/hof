---
title: "Modules"
description: "Dependency management for CUE and code generators"
brief: "Dependency management for CUE and code generators"

weight: 20
---

### Hof & CUE Modules

Every __hof generator__ is also a __CUE module__,
and in fact, many of hof's other features can
be used from the module system too.
Hof & CUE's modules serve the same purpose as other languages,
allowing to to version, share, and reuse code.
"`hof mod`" is the subcommand for working with modules and dependencies.
The implementation is based on Go modules.

The name of a module should be the same the git repository.
`hof` talks directly to git repositories and many of
`hof`'s commands will accept modules as an input argument too.

##### [To learn more, see the modules section](/modules/).

<br>

{{<codeInner>}}
# create a new module
hof mod init github.com/hofstadter-io/example

# add a dependency
hof mod get github.com/hofstadter-io/hof@v0.6.8
  or
hof mod get github.com/hofstadter-io/hof@latest

# tidy dependencies
hof mod tidy

# fetch dependencies
hof mod link
  or
hof mod vendor
{{</codeInner>}}

