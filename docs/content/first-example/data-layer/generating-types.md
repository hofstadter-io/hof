---
title: "Generating Types"
brief: "for the server"
weight: 20
---

To generate types, we need to do two things

- Update our generator definition
- Add the template for Go types

### Generator Changes

Add the following changes in their appropriate places into the existing generator definition.

{{< codePane title="gen/server.cue" file="code/first-example/data-layer/content/gen/type.html" >}}


### Type Template

The following creates

- a Go struct for our Model
- a Go map for storing instances of the type
- several Go functions as helpers for the data store

Create a new template called `type.go`

{{< codePane title="templates/type.go" file="code/first-example/data-layer/content/templates/type.go" lang="go" >}}

### Regenerate the server

You can now run `hof gen ./example` and you should find a new `./output/types` directory.

