---
title: "Code Generation"
linkTitle: "Code Gen"
brief: "Render any file from data and templates."
description: "Render any file from data and templates."
keywords:
  - text/template
  - template-based
  - ad-hoc rendering
weight: 10
---


{{<lead>}}
Code generation is central to `hof`.
`hof`'s ad-hoc code generation feature enables developers to
__generate files for any programming language or framework__
by rendering templates with data.
This page introduces code generation through `hof`.
{{</lead>}}

With `hof gen`, you can combine CUE, Yaml, or JSON arguments
with templates to generate files and directories of code automatically.
Through flags and configurations, you have complete control over
the file generation process and the final content created.

{{<codeInner lang="sh">}}
hof gen <cue, data, config>... [--flags...]
{{</codeInner>}}

`hof`'s template-based code generation has two approaches:
ad-hoc and a more sophisticated configuration or module-based method.

1. `hof gen -T` is code gen from flags
2. `hof gen -G` is code gen from config


In this section, we will explore the ad-hoc method which used the `-T` flag.
The [first example](/first-example/) is a step-by-step guide
on creating a generator using CUE based configuration.
The [getting-started/create](/getting-started/create/) section
will introduce `hof create`, which enables running generators directly from git repositories.

Each of these methods utilizes the same core concepts and processes, and they are suitable for different use cases:

- `hof gen -T` is ideal for simpler scenarios when just a few files are involved.
- `hof gen -G` is designed for large-scale code generation and reusable modules.
- `hof create` is intended for one-time setup and application bootstrapping.

[Code generation topics](/code-generation/) are discussed in a dedicated section.


## Data + Templates

By running the command `hof gen interlude.json -T interlude.template`,
users can perform ad-hoc template rendering to __combine any data source with any template__.

{{<codePane3
  title1="interlude.json" file1="code/getting-started/code-generation/interlude.json" lang1="json"
  title2= "interlude.template" file2= "code/getting-started/code-generation/interlude.template" lang2= "txt"
  title3="> terminal" file3="code/getting-started/code-generation/interlude.txt" lang3="txt"
>}}

#### `hof`'s templates are built on Go's `text/template` package with [extra helpers](/code-generation/template-writing/helpers/) added.

<br>

By using a hyphen symbol "`-`", you can stream any data into `hof gen`.
This feature can be helpful when you need to render an API response or process the output of a command.

{{<codeInner lang="sh" title="> terminal">}}
# use '-' to send output from another program to hof
$ curl api.com  | hof gen - -T template.txt

# intermix piped input with other entrypoints
$ cat data.json | hof gen - schema.cue -T template.txt

# set the data format when needed (cue help filetypes)
$ cat data.yaml | hof gen yaml: - schema.cue -T template.txt
{{</codeInner>}}


### Writing to file

Use  `=` (equal) to write to a file by name.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt=output.txt
{{</codeInner>}}

If you want to write all outputs to a directory, use the `-O` flag.
Files will have the same name as the template if not set individually.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt -O out/
{{</codeInner>}}

You can set the filename from data by using an inline template as the output name. 
Make sure you "wrap it in quotes".

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt="{{ .name }}.txt"
{{</codeInner>}}

You can combine these options to control the output directory and the filename. 

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
the data and schemas from the CUE entrypoints. (we'll see this below)

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt -T =debug.yaml -O out/
{{</codeInner>}}


### Watching for Changes

Use the -`w/--watch` flag to watch for changes and re-render output.
Think of the `-w/--watch` flag as a live-reload option that
monitors your code and automatically re-renders your output when changes are detected.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -T template.txt -T debug.yaml -O out/ --watch
{{</codeInner>}}

There are additional watch flags available if automatic detection misses files.

## On Using CUE

<br>

`hof`'s inputs are `cue`'s inputs, or "CUE entry points".

The inputs hold CUE values, which can be intermixed with your data to apply schemas, enrich the data, or transform before rendering.
When running commands, the CUE entrypoints are combined into a single CUE value.
If you want to place data files in a specific path, use the `@path.to.location` syntax.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.json@inputs.data schema.cue -T template.txt -T =debug.yaml -O out/
{{</codeInner>}}

The final data passed to a template must be concrete or fully specified. 
This means the value needs to be like JSON data before template rendering accepts them.
As you will see, hof provides flexibility and control for how the CUE values are selected, combined, and mapped to templates.

We keep parity with `cue`, so tooling from the broader ecosystem
still works on our inputs and reduces context-switching costs.
With hof, you can take advantage of all the possibilities and power of CUE
for selecting, combining, and mapping the values to templates for rendering.

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

You can control how input data and schemas are combined with templates
and written to files by using the flexible format of the `-T` flag for code generation.


### Selecting Values and Schemas

Use `:<path>` to select a value and `@<path>` to apply a schema

We can remove the `.Input` from our templates and pick the data and schema with flags.
This approach can be helpful when we don't have control over the input data or if it comes in a data format.

{{<codePane2
  title1= "types.go" file1= "code/getting-started/code-generation/typesInput.go" lang1= "go"
  title2="> terminal" file2= "code/getting-started/code-generation/out/typesInput.go" lang2= "text"
>}}


### Partial Templates

Partial templates are fragments that are used in other templates.
Unlike regular templates, these do not map to an output file. You can capture repeated sections
like the fields of a struct or the arguments to a function.


Additionally, partials can invoke other partials, which is helpful for modularizing templates into logical components.

There are two ways to define and use partial templates:

- Use the `{{ define "name" }}` syntax in a regular template
- User the `-P` to load them from a file

For example, we _could_ extract the generation of fields into its template.
We won't here, but as an example could include complex tasks like struct tags for Go fields. 

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen data.cue schema.cue -P field.go -T types.go -O out/
{{</codeInner>}}

{{<codePane3
  title1="types.go" file1="code/getting-started/code-generation/typeP.go" lang1="go"
  title2= "field.go" file2= "code/getting-started/code-generation/field.go" lang2= "go"
  title3= "out/types.go" file3= "code/getting-started/code-generation/out/types.go" lang3= "text"
>}}


### Repeated Templates

In addition to looping over data and applying a template fragment, you can use
__repeated templates__ to write a file for each element in a CUE or data iterable,
such as a list or struct field.
This is useful for creating a separate file for each item,
be it a type, API route, UI component, or DB migration.

To render a file for each element in the input to a `-T` flag, use `[]`.

{{<codeInner title="> terminal" lang="sh">}}
$ hof gen types.cue schema.cue -T type.go="[]{{ .Name }}.go" -O out/
{{</codeInner>}}

{{<codePane3
  title1= "type.go" file1= "code/getting-started/code-generation/type.go" lang1= "go"
  title2= "out/User.go" file2= "code/getting-started/code-generation/out/User.go" lang2= "go"
  title3= "out/Post.go" file3= "code/getting-started/code-generation/out/Post.go" lang3= "go"
>}}



### Understanding the -T flag

The `-T` flag in hof gen has a flexible format that
allows you to customize your templates' input data, schema, and output path. 

`-T "<template-path>:<input-path>[@schema-path]=<out-path>"`

This flag allows you to

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
[hof ad-hoc render tests](https://github.com/hofstadter-io/hof/tree/_dev/test/render).



## What are Generators and Modules


Generators are configurations for the `hof gen` flag, typically defined
in CUE modules or Git repositories. 

To turn your ad-hoc `hof gen ... -T ...` commands into a generator, 
add `--as-module <module name>` to the end of your current flag.


{{<codeInner title="> terminal">}}
$ hof gen ... --as-module github.com/username/foo
{{</codeInner>}}

Several files are generated, including a CUE file that houses your generator and additional files for configuring a CUE module.

{{<codePane file="code/getting-started/code-generation/adhoc-mod-snippet.html" title="generator.cue snippet">}}

The next page will provide an overview of modules in general,
and the [first-example](/first-example/) offers detailed instructions for creating a generator from scratch.

