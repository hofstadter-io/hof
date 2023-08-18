---
title: Data Modeling

weight: 40
---

{{<lead>}}
Data models are core to our applications and architectures.
With `hof`, you specify their types, shapes, and relations using CUE
and will combine, checkpoint, and track history
as your data models and applications evolve.
They become the source of truth when generating or updating
code across yor technology stack or service fleet.
{{</lead>}}

`hof datamodel` features enable:

- checkpoints, history tracking, structural diff calculations
- auto generated database migrations
- client/server version negotiation and request/response upgrades
- data lenses and other transformations
- More than CRUD generation from relations and history

## Base Datamodel

The base data model is any CUE Value with
`#hof: datamodel: root: true` set.

{{<codePane3
    title1="dm.cue (by hand)"      file1="code/data-modeling/by-hand.html"
    title2="dm.cue (by import)"    file2="code/data-modeling/by-import.html"
    title3="dm.cue (by attribute)" file3="code/data-modeling/by-attribute.html"
>}}

`hof` leaves those details to the user and module authors.
We also provide some common schemas and extras you can use
along with the main data model features `hof` provides.


## Extending Datamodels

With `hof datamodels` powered by CUE, you will be able to:

- enforce schemas and provide defaults
- share them as modules across applications
- merge or unify them like any CUE value
- enrich them conditionally before code gen

The first three points are native to CUE,
so you have all the same capabilities when using `hof`.
Modules are currently only available if `hof`.
CUE is working on their own dependency management
and we are involved in that work as well.
Learn more in the [modules section](/modules/).

When it comes to enriching, there are a few ways or places this happens.
`hof`'s base datamodel doesn't provide any structure,
in part because there is not one to rule them all,
but also because we want to let users define their own
while still getting the checkpointing, history, and diff features.

### Providing structure

More often than not, you will want to
provide structure to the datamodel.
[hof/schema/dm/sql](https://github.com/hofstadter-io/hof/blob/_dev/schema/dm/sql/dm.cue) is one example of this.
You define the structure by using the `@history()`
on a collection within your datamodel.
This is a tracking and pivot point.
`hof` will manage history and diffs
on a structural level, at each history point.


{{<codePane title="providing-structure.cue" file="code/data-modeling/providing-structure.html" >}}

`hof` will manage each CUE field in a value with `@history()`.
In the above code, this would be

- each model in `Models`, like `User`
- each field in `User.Fields`
- each view in `Views`

### Enriching values

Generally, the input `Datamodel` to a `Generator` will use generic field types.
When generating code, it can be helpful to enrich these values to calculate
template output in CUE rather than the templates themselves.
You can apply these by wrapping the datamodel or applying them in your generator inputs.

Places where this is helpful are:

- mapping types, especially collections and relations, to the target language or technology.
- adding various string casings or manipulations

An example of this is mapping `hof`s `schema/dm/fields` like `uuid`
to a package in languages like Go or Python

{{<codePane3
    title1="schema/dm/enrichers/go.cue" file1="code/hof-schemas/dm/enrichers/go.html"
    title2="schema/dm/enrichers/py.cue" file2="code/hof-schemas/dm/enrichers/py.html"
    title3="your-module/schema/dm.cue"  file3="code/data-modeling/using-enrichers.html"
>}}


### Future ideas

We plan to make data models usable in other `hof` sub-systems

- more features around multiple datamodels. Multiple are supported today, we could provide mappings from an instance in one to the other.
- `hof/flow` to enable similar transparent, inter-version transforms
- `hof/chat` to provide generator specific schemas for LLMs to target (like Microsoft/Guidance)

## Datamodel Commands

{{<codeInner title="example command usage" lang="shell">}}
$ hof dm list   (print known data models)
NAME         TYPE       VERSION  STATUS  ID
Config       object     -        ok      Config
MyDatamodel  datamodel  -        ok      datamodel-abc123

$ hof dm tree   (print the structure of the datamodels)

$ hof dm diff   (prints a tree based diff of the datamodel)

$ hof dm checkpoint -m "a message about this checkpoint"

$ hof dm log    (prints the log of changes from latest to oldest)

You can also use the -d & -e flags to subselect datamodels and nested values
{{</codeInner>}}

## Datamodels and Code Generation

Datamodels form the basis for input to code generation.
Typically, a generator will require a `dm.Datamodel` and
some other inputs specific to that generator.
You can then use them in the templates like any other value.

<!--The [checkpointing & history](/data-modeling/checkpointing-and-history/) page will cover using these during code gen.-->

Here is an example snippet that you would use with our
[supacode generator for full stack applications](https://github.com/hofstadter-io/supacode).

{{<codePane title="using the supacode generator" file="code/data-modeling/dm-and-code-gen.html" >}}

