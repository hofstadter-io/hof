---
title: "Define (impl)"
linkTitle: "Define (impl)"
weight: 3
---


Generators start by creating a _generator definition_.

1. Start a Cue module
2. Start a generator definition
3. The schema for a generator
4. Filling out the generator


### 1. Start a Cue module

Modules work like Go modules and have similar rules.

```
hof mod init cue domain.com/example/server
```

[cue mod file]

```
hof mod vendor cue
```

[more info on modules]


### 2. Start a generator definition

{{<codePane title="full-example/server/gen/gen.cue" file="code/getting-started/full-example/tmp/gen-def.html" >}}

### 3. The schema for a generator

Generators have a `schema` to define a contract between
the implementor and the user.
We will define the schema for our server in the next section.

Similarly, __hof__ has a `schema` for generators
to define the contract with implementors.
When you write a generator, it must
implement this schema.

Brief and link to these:

[HofGenerator schema]

[HofGeneratorFile schema]

### 4. Filling out the generator
