---
title: "The Data Layer"
brief: "managing your data model"

weight: 20
---

{{<lead>}}
Data modeling is core to the development process.
As our understanding of the design evolves,
so must our code.
{{</lead>}}

{{<lead>}}
`hof` has a data modeling system that works with the code generation process.
{{</lead>}}


- define data models in CUE, couple these with and into code generators
- generate language types, database tables, libraries, API handlers, and more
- checkpoint the data model and maintain a history for version transforms and migrations

This section expands on our `simple-server` to use `hof/dm.#Datamodel`.
We will first

1. Create a todo application data model
1. Generate Go types and a simple Library
1. Start with a Go map for storage, later a database
1. Create CRUD routes for the datamodel

After we will see `hof`'s code _regeneration_ capabilities by

1. Customizing the generated code
1. Updating the data model
1. Regenerating our application 
1. See how `hof` fits into typical application development

Finally, we will look at how to upgrade our generator to use a database.
Automatic migrations are covered in the [model history section](/first-example/model-history/).

The full code for this section can be found on GitHub
[code/first-example/data-layer](https://github.com/hofstadter-io/hof-docs/tree/main/code/first-example/data-layer)

_Database storage and automatic CRUD handler generation
will be covered in more advanced sections._

{{<childpages childBriefs="true">}}

