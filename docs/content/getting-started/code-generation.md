---
title: "Code Generation"
linkTitle: "Code Gen"
description: "Render any file from data and templates."
brief: "Render any file from data and templates."

weight: 10
---


{{<lead>}}
Code generation is central to `hof`.
We can generate any file for any language or framework by rendering templates with data.
{{</lead>}}

`hof gen` accepts CUE, Yaml, and JSON arguments
and combines them with templates to generate files and directories of code.
This process is highly flexible and controlled through flags and configuration,
or what we call generators.

{{<codeInner lang="sh">}}
hof gen <cue, data, config>... [--flags...]
{{</codeInner>}}


There are two main modes for
template-based code generation.

1. `hof gen -T` is code gen from flags
2. `hof gen -G` is code gen from config

We will look at `hof gen -T` in this section.
The [first example](/first-example/) is a walkthrough
in writing a generator (the `-G` flag).
The [getting-started/create] section will introduce
`hof create` for running generators directly from git repositories.

The concepts and processing are the same across all of them and
each has a use case it is best at:

- `hof gen -T` is suitable for simple cases or when you don't want dependencies
- `hof gen -G` is aimed at reusable and sharable modules with dependencies
- `hof create` is intended for interactive one-time setup and bootstrapping

[Code generation topics](/code-generation/) are discussed in a dedicated section.


## Data + Templates

`hof gen interlude.json -T interlude.template` is ad-hoc template rendering.
Combine any data source with any template.

{{<codePane3
  title1="interlude.json" file1="code/getting-started/code-generation/interlude.json" lang1="json"
  title2= "interlude.template" file2= "code/getting-started/code-generation/interlude.template" lang2= "txt"
  title3="> terminal" file3="code/getting-started/code-generation/interlude.txt" lang3="txt"
>}}

#### `hof`'s templates are built on Go's `text/template` package with [extra helpers](/code-generation/template-writing/) added.

<br>

You can also pipe any data into `hof gen` by using a "`-`" (hyphen).
This can be helpful when you want to render
an API response or process command output.

{{<codeInner lang="sh" title="> terminal">}}
# use '-' to send output from another program to hof
$ curl api.com  | hof gen - -T template.txt

# intermix piped input with other entrypoints
$ cat data.json | hof gen - schema.cue -T template.txt

# set the data format when needed (cue help filetypes)
$ cat data.yaml | hof gen yaml: - schema.cue -T template.txt
{{</codeInner>}}


### Writing to file

Use  `=` (equal) after the template name to write to a file.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt=output.txt
{{</codeInner>}}

Use `-O` to write all outputs to a directory.
Files will have the same name as the template if not set individually.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt -O out/
{{</codeInner>}}

The output name can be a template so that you can control the filename from the data.
Make sure you "wrap it in quotes".

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt="{{ .name }}.txt"
{{</codeInner>}}

These can be combined so you can control
where the output goes and how files are named.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue \
  -O out \
  -T template.txt="{{ .name }}.txt"
{{</codeInner>}}

### Write Data Files

Omit the template path at the beginning
and `hof` will infer the data format
from the output file extension

{{<codeInner title="> terminal" lang="sh">}}
# full value to a single data file
$ hof gen data.cue schema.cue -T =data.yaml

# data file per item in iterable value
$ hof gen data.cue schema.cue \
  -O out \
  -T :items="[]{{ .name }}.json"
{{</codeInner>}}

### Multiple Templates

You can use the `-T` flag multiple times.
Each is independent and can have different options for 
the data and schemas from the CUE entry points. (we'll see this below)

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt -T =debug.yaml -O out/
{{</codeInner>}}


### Watching for Changes

Use `-w`/`--watch` to observe for changes and re-render output.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt -T debug.yaml -O out/ --watch
{{</codeInner>}}

There are extra watch flags if automatic detection doesn't entirely work.


## On Using CUE

<br>

`hof`'s inputs are `cue`'s inputs, or "CUE entry points".

The inputs hold CUE values, which can be intermixed with your data to
apply schemas, enrich the data, or transform before rendering.
When running commands, the CUE entry points are combined into one considerable CUE value.
The final data passed to a template must be concrete or fully specified.
This means the value needs to be like JSON data before template rendering accepts them.
As you will see, `hof` provides you flexibility and control
for how the CUE values are selected, combined, and mapped to templates.

We keep parity with `cue` so tooling from the broader ecosystem
still works on our inputs and reduces context-switching costs.

__You can safely use all the possibilities and power of CUE here.__



## Setup for Examples

<br>

We will be using the following inputs in the examples below.
We define __a schema__ and write our types as __data values__ in CUE.

#### Schema & Data

{{<codePane2
  title1= "schema.cue" file1= "code/getting-started/code-generation/schema.html"
  title2= "data.cue" file2= "code/getting-started/code-generation/data.html"
>}}

We can use `cue` to see what the full data looks like


<details>
<summary><b>$ cue export data.cue schema.cue</b></summary>
{{<codePane title="> terminal" file="code/getting-started/code-generation/full-data.json" lang="text">}}
</details>

<br>

#### Starting Template

{{<codePane2
  title1= "types.go" file1= "code/getting-started/code-generation/types.go" lang1= "go"
  title2="> terminal" file2= "code/getting-started/code-generation/out/types.go" lang2= "text"
>}}


## Controlling Code Generation

The `-T` flag has a flexible format so you can
control how the input data and schemas are
joined with templates and written to files



### Selecting Values and Schemas

Use `:<path>` to select a value and `@<path>` to apply a schema

We can remove the `.Input` from our templates and
pick the data and schema with flags.
This is helpful if we do not control the input data
or if it comes in a data format.

{{<codePane2
  title1= "types.go" file1= "code/getting-started/code-generation/typesInput.go" lang1= "go"
  title2="> terminal" file2= "code/getting-started/code-generation/out/typesInput.go" lang2= "text"
>}}


### Partial Templates

Partial templates are fragments
that are used in other templates.
Unlike regular templates, these do not map to an output file. You can capture repeated sections
like the fields to a struct or the arguments to a function.

Partials can also invoke other partials,
which makes them ideal for breaking up
your templates into logical components.

There are two ways to define and use partial templates:

- Use the `{{ define "name" }}` syntax in a regular template
- User the `-P` to load them from a file

Let's extract field generation into its template, where we _could_ make it complex.
We won't here, but an example is struct tags for our Go fields.
We can also use template helpers in the output file path.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -P field.go -T types.go -O out/
{{</codeInner>}}

{{<codePane3
  title1="types.go" file1="code/getting-started/code-generation/typeP.go" lang1="go"
  title2= "field.go" file2= "code/getting-started/code-generation/field.go" lang2= "go"
  title3= "out/types.go" file3= "code/getting-started/code-generation/out/types.go" lang3= "text"
>}}


### Repeated Templates

We just saw how to loop over data and apply a template fragment.
We can also render with __repeated templates__, which are processed
for each element of an iterable (list or struct fields)
and write a file for each component.

Use `[]` to render a file for each element in the input to a `-T` flag.


{{<codeInner title="> terminal" lang="sh">}}
$ hof gen types.cue schema.cue -T type.go="[]{{ .Name }}.go" -O out/
{{</codeInner>}}

{{<codePane3
  title1= "type.go" file1= "code/getting-started/code-generation/type.go" lang1= "go"
  title2= "out/User.go" file2= "code/getting-started/code-generation/out/User.go" lang2= "go"
  title3= "out/Post.go" file3= "code/getting-started/code-generation/out/Post.go" lang3= "go"
>}}



### -T Flag Details

The `-T` flag for `hof gen` has a flexible format:

`-T "<template-path>:<input-path>[@schema-path];<out-path>"`

This flag enables you to

1. Render multiple templates by using `-T` more than once
1. Select a value with `:<input-path>`
1. Select a schema with `@<schema-path>`
1. Write to file with `=<out-path>`
1. Control the output filename with `="{{ .name }}.txt"`
1. Render a single template multiple times with `="[]{{ .filepath }}"`

<br>

{{<codeInner title="-T variations" lang="txt">}}
hof gen input.cue ...

  # Generate multiple templates at once
  -T templateA.txt -T templateB.txt

  # Select a sub-input value by CUEpath (. for root)
  -T templateA.txt:foo
  -T templateB.txt:sub.val

  # Choose a schema with @
  -T templateA.txt:foo@#foo
  -T templateB.txt:sub.val@schemas.val

  # Writing to file with =
  -T templateA.txt=a.txt
  -T templateB.txt:sub.val@schema=b.txt

  # Templated output path 
  -T templateA.txt='{{ .name | lower }}.txt'

  # Repeated templates are used when
  # 1. the output has a '[]' prefix
  # 2. the input is a list or array
  #   The template will be processed per entry
  #   This also requires using a templated out path
  -T template.txt:items='[]out/{{ .filepath }}.txt'
{{</codeInner>}}

You can find more examples in the
[hof render tests](https://github.com/hofstadter-io/hof/tree/_dev/test/render).



## What are Generators and Modules

Generators are `hof gen` flags as configuration,
often in CUE modules and git repositories.
The next page will overview modules more generally.
The [first-example](/first-example/) details
creating a generator from scratch.

To turn your ad-hoc `hof gen ... -T ...` commands into a generator
by adding `--as-module <module name>` after the current flags.

{{<codeInner title="> terminal">}}
$ hof gen ... --as-module github.com/username/foo
{{</codeInner>}}

You will see a few files created.
There will be a CUE file that contains your generator
and a few others for setting up a CUE module.

{{<codePane file="code/getting-started/code-generation/adhoc-mod-snippet.html" title="generator.cue snippet">}}
