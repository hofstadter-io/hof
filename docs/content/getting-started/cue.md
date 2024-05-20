---
title: "CUE"
description: "A declarative language for configuration and code generation."
brief: "A declarative language for configuration and code generation."
keywords:
  - declarative code generation 
  - cue 
  - generators
  - declarative schemas 
  - reusable modules
  - configuration-style syntax
  - dependency management
weight: 5
---

{{<lead>}}
`hof` leverages the power of CUE (Configure, Unify, Execute) to define inputs, schemas, and generators.
CUE is a language specifically designed for configuration and is an integral component of the `hof` developer experience.
{{</lead>}}

Two key considerations drove the choice to use CUE:

1. __Declarative Code Generation__<br>
   `hof` uses code generators defined in CUE, combining data and templates to generate files.
   With CUE as the source of truth, we gain consistency for inputs and mappings, allowing
    composable generators to generate code across languages and technologies.
    `hof` allows you to define the declarative schemas to code generators and package them into reusable modules.

2. __Purpose-Built for Large-Scale Configuration__<br>
   CUE has a strong foundation and robust tooling for managing configuration across multiple files, packages, and modules.
   CUE's logical nature ensures that the configuration is both efficient and correct.
    Since `hof`'s inputs are typically declarative configurations, CUE is a natural fit throughout the process.

Other benefits of CUE include:

- Configuration-style syntax that is familiar, yet fresh.
- Built on sound principles and rich heritage.
- Declarative and logical, providing confidence in the validity and consistency of the configuration.
- Purpose-built language and tooling for large-scale configuration.
- Dependency management for data, schemas, configuration, and code generation.

For those unfamiliar with CUE, Bitfield Consulting has a great
[Introduction to CUE from a JSON perspective](https://bitfieldconsulting.com/golang/cuelang-exciting).
It is a quick read and provides an excellent foundation for using CUE with `hof`.

To further your CUE knowledge, be sure to check out these resources:

- [CUE documentation](https://cuelang.org) (from the CUE Team)
- [Cuetorials](https://cuetorials.com) (by Hofstadter)


### Hof's CUE commands

Hof embeds CUE's `vet, def, eval, export` commands for your convenience.
The are mostly drop in alternatives, some codecs are not available and
several enhancements have been added.

The enhancements are:

- additional methods for data placement
- increased flexibility for environment variables 
- @userfiles() to include any file
- `--tui` flag to open hof's TUI for the commands


### Hof & CUE Modules

__hof__ has a preview version for __CUE modules__.
Hof & CUE's modules serve the same purpose as other languages,
allowing to to version, share, and reuse code.
CUE's module system is still largely experimental.
We will eventually migrate over once sufficient features
are in place, at which point we will provide automation to update.

Most of hof's features can be used from the module system.
"`hof mod`" is the subcommand for working with modules and dependencies.
The implementation is based on Go modules.

The name of a module should be the same the git repository.
`hof` talks directly to git repositories and many of
`hof`'s commands will accept modules as an input argument too.
There is also support for OCI based repositories.

##### [To learn more, see the modules section](/modules/).

<br>

{{<codeInner>}}
# create a new module
hof mod init github.com/hofstadter-io/example

# add a dependency
hof mod get github.com/hofstadter-io/hof@v0.6.8
  or
hof mod get github.com/hofstadter-io/hof@latest

# tidy dependencies
hof mod tidy

# fetch dependencies
hof mod link
  or
hof mod vendor
{{</codeInner>}}
