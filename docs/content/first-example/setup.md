---
title: "Setup"
brief: "a new hof generator"
weight: 5
---

{{<lead>}}
Generators are made of some CUE files and directories of templates.
You can run `hof gen init <name>` to bootstrap a new generator module.
{{</lead>}}

### Bootstrapping Your Geneartor

Throughout the first example, you will be
building a code generation module.
The following commands will create a new
directory, git repository, and hof generator
for you to start with.

{{<codeInner title="Setup Commands">}}
// create a directory
mkdir example && cd example
git init   // (optional)

// initialize a generator module
// the name is often the same as a github repo
hof gen init hof.io/docs/example

// temporary fix until v0.6.8
mv generators gen
{{</codeInner>}}

Your working directory should now look like:

{{<codeInner title="Module layout">}}
example/
|  // default directories
├── creators/
├── examples/
├── gen/
├── partials/
├── schema/
├── statics/
├── templates/
|
|  // dependency files
├── cue.mods
├── cue.sums
└── cue.mod/
{{</codeInner>}}

You will often have other files depending on the languages or technologies you choose.
A common example are the files for dependency management like `package.json` and `go.mod`.
There are no restrictions or limits on what you can include.

The default directories:

- `schema` is where your schemas for the other parts go
- `gen` holds CUE files for specifying your generators
- `templates`, `partials`, and `statics` are files for generators
- `creators` are often in a subdirectory with their own templates
- `examples` for using your generators, also helpful for testing
