---
title: "Generator"
brief: "for mapping In to Out"
weight: 15
---

{{<lead>}}
Generators specify the mapping from inputs through templates to output files.
Using the schema and user input, you decide what files to generate,
what inputs to provide to templates, and can reshape or enrich the inputs.
Your generator file definds this process to
map the `In` field to the `Out` field for `hof` to process.
{{</lead>}}

<br>

### Generator Schema

HofGenerators and Files have their own schemas.
Only the core fields are shown here.
You can find the real schemas in [the hof repository](https://github.com/hofstadter-io/hof/tree/_dev/schema/gen).

{{<codePane file="code/first-example/simple-server/content/generator/schemas.html" title="Generator Schemas">}}

### Mapping In to Out

Your goal as a generator writer is to fill in the `Out` field from the `In` field.
The `In` field will be passed through the templates and partials
if they are defined in the `Out` list.

##### `In` represents the data presented to the templates. 

- You can have your users set values in meaningful fields and collect them
- You can add any other separate or calculated data
- The top-level `In` is passed to all templates and can be extended with the per-template field

##### `Out` holds the files to be rendered by hof

- Is a list of `HofGeneratorFile`s, where each...
- Has a filepath for where it will be written under `Outdir`
- Has an `In` value which will be unified with the top-level

There are typically two types of files (templates)

- Files which are generated once per application. Think of a `main.go` or `index.js`. We often call these "once" files or templates.
- Files which are generated for each sub-resource. These might be routes in a server or commands in a CLI. We often call these "repeated" files or templates.

##### You can use as many helper fields or calculations as you want

- Calculate commonly used fields or values for `In`
- Interpolate common base paths for filepaths
- Create file lists for repeated templates.

### Server Generator

The following is the generator for our simple REST server.

{{<codePane file="code/first-example/simple-server/gen/server.html" title="gen/server.cue">}}


