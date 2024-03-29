---
title: "Setup"
brief: "a new generator"
weight: 5
---

{{<lead>}}
Generators are made of some CUE files and directories of templates.
You can run `hof gen init <name>` to bootstrap a new generator module.
{{</lead>}}

### Bootstrapping Your Generator

Throughout the first example, you will be
building a code generation module.
The following commands will create a new
directory, git repository, and hof generator
for you to start with.

{{<codeInner title="Create working directory">}}
mkdir example && cd example
git init      // (optional)
{{</codeInner>}}

{{<codeInner title="Initialize a generator module">}}
// the name is often the same as a github repo
hof gen init hof.io/docs/example
{{</codeInner>}}

Your working directory should now look like this:

{{<codeInner title="Module layout">}}
example/
|               // default directories
├── creators/
├── examples/
├── gen/
├── partials/
├── schema/
├── statics/
├── templates/
|
└── cue.mod/    // CUE dependencies
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
