---
title: Module Layout

weight: 10
---

{{<lead>}}
CUE modules have a specific set of files and directories
that make up the layout.
{{</lead>}}

### cue.mod

Both `cue` and `hof` will look for this directory and file,
by walking up the directory tree until it finds them.

The `cue.mod` directory has a few important subdirectories

{{<codeInner title="cue.mod directory">}}
cue.mod/
	module.cue  // indicates a CUE module
	sums.cue    // hashes for integrity
  /pkg        // where dependencies go
	/gen        // where `cue import` puts code
	/usr        // for extending or correcting /gen
{{</codeInner>}}

### module.cue

The `module.cue` file is how `cue` knows the current module name
and how `hof` records dependencies.

{{<codePane file="code/modules/cue.mod/module.html" title="cue.mod/module.cue">}}



### packages

A CUE module is a collection of packages.
These are similar to Go packages, and if
you use them the same, then they will behave the same.
However, CUE modules have several extensions.
We recommend sticking to the Go style packages and imports,
as this style is much easier for non-experts to understand.
If you want to learn about the other variations,
check out our page on [cuetorials.com - modules & packages](https://cuetorials.com/first-steps/modules-and-packages/).


The Go style of packages and imports:

- use lowercase and underscores
- have only one package per directory
- name should be the same as the directory
- use `<module-path>/<package-path>` for all imports

For example, in the `github.com/hofstadter-io/hof` module,
we have imports like:

{{<codePane file="code/modules/hof.html" title="github.com/hofstadter-io/hof - hof.cue">}}

