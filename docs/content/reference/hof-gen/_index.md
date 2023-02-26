---
title: "hof / gen"
linkTitle: "hof / gen"
weight: 30
---

Code generation is the core reason __hof__ was created.
It's how we turn single-source of truth designs into all the things.
It was needed so that when our datamodel changes,
we could edit just one or two design files
rather than dozens of source code, test files, and database schemas.

We are working on this concurrently with the datamodels work.
Like that section, we will post links here as they become available.

The goal here is to emit code from Cue in a general and language agnostic way.

This is the best example right now: https://github.com/hofstadter-io/hofmod-cli (edited) 

These generators basically use the Go text/template package behind the scenes

Schemas can be found in [./hof/schema/...](https://github.com/hofstadter-io/hof/tree/_dev/schema)

Using a DSL on top of Cue will enable more context, and thus capabilities, for code generation.
A DSL is equivalent to the Schema found in the generators, which is the data input to the template rendering

My intention for the datamodel work is to create a core which can be extended with context by domain schemas and fed into their generators. Another goal is to enable new code generators without modifying hof's code (the binary that you run)
This way, I can have a single data model and feed it into several generators, so that it is a "source of truth" to many outputs (i.e. client, server, sql)


Using Cue to specify the schema (DSL), introspect the user's design (using a schema / generator), and deciding what files should be created (gen dirs and templates). The actual process of pumping the data through the templates is in Go (lib/{gen.go,runtime.go,gen/})

There's also a shadow directory, so that we can 3-way diff merge between:

1. previous code
2. next code
3. custom code in the generated output

It was key for me to be able to write code in the output, change the design, regenerate, and have all things be OK `*`

`*` changing datamodels may require you to change custom code using it

There are edge cases around updating the generator and designs at the same time, would like to detect that situation and complain to the user more.

Renaming certain things is hard too (while keeping custom code) i.e. a file name can change, which will create a new file and delete the old one. Git is your friend :]

designs/ are filled in schema/ files which are the inputs to gen/ files which do some Cue magic to setup a data struct in hof, which itself just runs the Out field of generators through a process. Mostly Cue logic up front and the text/template logic in the templates/ and partials/


{{< childpages >}}
