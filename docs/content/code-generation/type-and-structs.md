---
title: "Type and Structs"
description: "Render types in multiple languages from CUE with hof"
brief: "with CUE + text/template"

weight: 80
---

{{<lead>}}
One of the common questions for CUE is
how to generate the matching types in a given language.
We will introduce the ideas and complexities with type generation
while also showing concrete examples of `hof gen -T` variations.
{{</lead>}}

Types are central to the languages we program in.
`hof` enables us to specify the shape and rules for
our types in CUE and then generate them in one or more languages.
The schemas and generators for types are the foundation
for many other generators and are important to understand.

## CUE Discussions

More has been written and discussed at the following links.
Keep these handy as you read through.

- https://github.com/cue-lang/cue/discussions/1865
- https://github.com/cue-lang/cue/discussions/1889

## Generating Types

CUE is a great language for defining types and their validations,
but how do we turn them into the structs, classes, and other
language specific representations?

Right now, and generally, the answer is `text/template`.
Right new because CUE does not have this capability.
Generally because CUE cannot capture
the variety and nuances between languages.
What are these complications?

- OOP vs structural typing, how do you express inheriance or embedding?
- CUE's types often look like a Venn Diagram with a languages types
- Native validation will be faster, will also need to be generated.
- Casing preferences per language
- Public, private, and protected visibility
- Default values, when and where they are setup

It would be a burden to put this all on CUE developers to figure out and maintain.
By using text interpolation, we can generate types without modifying CUE.
_Note, CUE does intend to support some language targets, but there
is no timeline for when this will happen yet or what it will look like._

If we want to have a single-source of truth, we need two things

1. An abstract representation, DSLs are a natural fit
1. Mappings to our target languages and technologies

CUE happens to be a good language and model for
writing and validating both the representation and mappings.


## Type DSLs

We believe that using a DSL, rather than native CUE expressions,
is the better option. There are many things which we cannot
express directly in CUE types and constraints, and using
attributes requires the tool to understand these.
So in order to provide maximal flexibility to experiment
without needing to modify `cue` or `hof`, we use DSLs.
Fortunatedly, CUE makes it easy to create and validate DSLs,
it's just a perspective of CUE values afterall.

Another hard question is "is there a single type schema or DSL to rule them all?"
Probably not, though one might be able to capture the majority of cases.
As defined, the type DSLs and schemas can be extended or specialized, like any CUE value.
This will give the community a way to combine and specialize them as needed.


### A Type Schema

With `hof`, we are building some reusable data model schemas.
This subsection will show you a simplified version for demonstration.

- schema
- example types used

{{<codePane title="A Type Schema" file="code/getting-started/code-generation/schema.html">}}

### Example Types

Let's use a blogging site as our example.

{{<codePane title="types.cue" file="code/getting-started/code-generation/data.html">}}

Run `cue eval types.cue schema.cue --out yaml` to see it's final form

{{<codePane title="types.cue" file="code/getting-started/code-generation/data.yaml" lang="yaml">}}


## The Templates

Now we have to implement the above schema
in our target languages and technologies.

We will run all of the following with `hof gen types.cue schema.cue -T ...`

Output will be put into the `out/` directory.


### Go Structs

We can start with a single template and file for all types.

Run `hof gen types.cue schema.cue -T types.go`

or `hof gen ... -T "types.go;out/types.go"` to write to file

{{<codePane2
	title1="types.go" file1="code/getting-started/code-generation/types.go" lang1="go"
	title2="out/types.go" file2="code/getting-started/code-generation/out/types.go" lang2="go"
>}}

We can render with __repeated templates__, which are processed
for each element of an iterable (list or struct fields).

Run `hof gen types.cue schema.cue -T "type.go;[]out/{{.Name}}.go"`

{{<codePane3
	title1="type.go" file1="code/getting-started/code-generation/type.go" lang1="go"
	title2="out/User.go" file2="code/getting-started/code-generation/out/User.go" lang2="go"
	title3="out/Post.go" file3="code/getting-started/code-generation/out/Post.go" lang3="go"
>}}

Use __partial templates__ for repetitious template content within a single file.
Let's extract field generation into its own template, where we _could_ make it complex.
We won't here, but an example is struct tags for our Go fields.
We can also use template helpers in the output filepath.

Run `hof gen types.cue schema.cue -P field.go -T "typeP.go;[]out/{{ lower .Name }}.go"`

{{<codePane3
	title1="typeP.go" file1="code/getting-started/code-generation/typeP.go" lang1="go"
	title2="field.go" file2="code/getting-started/code-generation/field.go" lang2="go"
	title3="out/user.go" file3="code/getting-started/code-generation/out/User.go" lang3="go"
>}}

### SQL & TypeScript

- multiple templates
- non-cue type ID (uuid, etc...)


### Protobuf

Show issue with indexing, consistent ordering

2 options

1. 


### More than types

1. REST & DB lib stubs (not just types)

- partials, introduce here, or earlier and expand here

### Generator Module

Show how to convert to a generator module

---



More advanced walkthrough and discussion in...



Briefly mention and link to

1. Generating types from more vanilla CUE (field: string, rather than DSL)
1. Generate for a framework

