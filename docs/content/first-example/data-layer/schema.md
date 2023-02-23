---
title: "Schema"
brief: "for the Data Model"

weight: 10
---

{{<lead>}}
Like generators, data models have a schema.
They capture the fields, relations, and other information
needed to generate types, database tables, and much more.
{{</lead>}}

### hof/dm.#Datamodel Schema

`hof`'s core `#Datamodel` schema is intentionally minimal.
It defines a hiearchy from `Datamodel > Model > Fields & Relations`.
The schema is sparce and open, having just what is needed
for `hof dm` to checkpoint and introspect data models.

The following as the core of the `#Datamodel` schema.

{{< codePane title="Core of #Datamodel" file="code/first-example/data-layer/content/schema/hof-dm.html">}}

See [hof datamodel schemas](/reference/hof-datamodel/schemas/) for the full schema.

### example/schema.#Datamodel

`hof`'s core `#Datamodel` schema is intended to be extended.
Since we will be using a Go map for a simple data store,
we will add a CUE field to `#Model` to track which
field to use for the index in our types.

We add a new file in our schema directory

{{< codePane title="example/schema/dm.cue" file="code/first-example/data-layer/content/schema/dm-type.html">}}

We now have an extended schema we can import to define our data models.

