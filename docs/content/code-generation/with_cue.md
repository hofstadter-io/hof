---
title: With CUE

weight: 100
---

Code Gen with CUE

- There are several types of "code gen" in relation to CUE
- data can use export with --out <data> or --outfile <file>.<ext>
- arbitrary uses same text/template
- cue export -e content --out text > file.txt
- write with tool/os

Missing:

- formatting
- diff3 & custom code
- expanded template helper functions
- shared partial templates, can still use partials defined within the template
- always writes, important when you use hot reload or make that checks FS events to compile


---
---

{{<lead>}}
There are several types of "code gen" in relation to CUE
{{</lead>}}

We should talk about these before focusing on `hof`.

Configuration or data (yaml,json)

- __Export__: CUE -> data (cue export)
- __Import__: data -> CUE (cue import)

Language types:

- __Get Go__: Go -> CUE   (cue get go)
- __TypeGen__: CUE -> Go   (... custom ...)

_some notes:
(1) CUE only has support for importing Go right now
(2) CUE only imports types from Go
(3) CUE does not have functions, so we cannot represent those in CUE without a DSL_

{{<lead>}}
Translating CUE to <lang> type (class or struct) is a challenging problem generally.
{{</lead>}}

It really can't capture all of the nuances in the vanilla form. We need to do something more complex.

### The Hof Method

Some helpful links

- Discussion on Slack about the data modeling prototype in the #general channel around here: https://hofstadter-io.slack.com/archives/C013WKK9W1F/p1640891812004200
- "Adding a Datamodel" section: https://docs.hofstadter.io/first-example/data-layer/
- https://docs.hofstadter.io/first-example/data-layer/relations/#commentary
- https://github.com/hofstadter-io/hofmod-server is based on the previous prototype of data modeling. It will be rewritten and broken into several modules
- A local code for `hofmod-sql` and `hofmod-gotype` are in the works. `hofmod-server` will then build on these.


Once `first-example` is done, several production versions of the concepts therein will be made into hofmods

### With CUE

There are also ways to generate types without hof

- one could do essentially the same thing hof is doing in pure CUE, though less sophisticated
- If you are willing to write Go, you can walk the AST and/or cue.Value to do some things. Attributes would make this more interesting.
- hof leaves more to the data model schema + generator, so it can be more flexible without modifying the tool. This should also be composable. I plan to have an example of that in the first-example/using-a-database or maybe another section


There are also discussion and issues on CUE GitHub.

- https://github.com/cue-lang/cue/discussions/1027
- https://github.com/cue-lang/cue/discussions/1038
- https://github.com/cue-lang/cue/discussions/482
- https://github.com/cue-lang/cue/issues/6
- https://github.com/cue-lang/cue/issues/943#issuecomment-1006480535


TODO, respond to

- https://github.com/cue-lang/cue/discussions/1167


