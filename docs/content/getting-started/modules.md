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

"`hof mod`" is the subcommand based on Go modules
for working with CUE modules and managing dependencies.
The logic and format is the same, with much code shared between the implementations.
Currently, while CUE is module aware and supports imports,
it does not yet have dependency management, but it will work the same as well.
This page has a brief overview. The commands will
be repeated when you need to run them on other pages.

The name of a module must be the same the git repository.
`hof` talks directly to git repositories and many of
`hof`'s commands will accept repositories as input too.

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

