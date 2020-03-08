# hof-lang

hof-lang is a framework for creating
low-code libraries, modules, and tools
for any technology, system, process, or platform.
It is modeled after
[Cuelang](https://cuelang.org)
with the unification to ensure correctness
and all the mathemtical rigor.
hof-lang differs in that it
makes code generation a first-class concept,
outside of the tool itself.
This means you can create domain specific,
low-code modules and tools for any language or technology
without needing to understand and modify the hof codebase.

hof-lang also has roots in
declarative programming and DSLs,
where your write your high-level
"designs" in the "language" of a DSL.
The __hof__ tool will read your designs
and use polyglot generators to output the implementation
in one or more computer languages and technologies.
You can modify the code in the output,
update your designs, and regenrate the output.
hof-lang allows you to work in both sides 
of the transformation, using diff3 to
merge the files together.


#### Concept flow in hof-lang

Designs -> DSLs -> Generators

__Designs__ are written in DSLs using __hof-lang__ syntax.
They form the "source-of-truth" for your
application, data validation, or other creation.
Designs are essentially data or configuration
with a predefined structure and their own validation from the DSL.

__DSLs__ have __hof-lang__ which acts as the spec for a domain or technology.
This DSL spec creates the contract between the designs and generators.
There are DSLs for data validation and creation, cross-platform CLIs,
REST APIs, CI/CD setups for projects.

__Generators__ implement a DSL in one or more languages or technologies.
They are the realization of a DSL from your designs.

hof-lang also supports modules, packages, and imports.
The semantics are based off of Golang, except exporting
is an explicit statement because data keys often start with lowercase letters.


#### Getting started

##### 1. Install the hof tool with

```bash
go get github.com/hofstadter-io/hof
```

##### 2. [Read over the documentation](./docs)
##### 3. [Explore the examples](./examples)
##### 4. [Join the conversation on Gitter](https://gitter.im/hofstadter-io/hof)
