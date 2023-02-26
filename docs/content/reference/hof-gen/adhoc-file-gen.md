---
title: "Adhoc File Gen"
linkTitle: "Adhoc File Gen"
weight: 25
---

{{<lead>}}
`hof gen` joins CUE with Go's text/template system and diff3.
This section focuses on the adhoc, one-liners
you can write to generate any file from any data.
{{</lead>}}

<br>

#### Learn about writing templates, with extra functions and helpers

[Template writing docs](/code-generation/template-writing/)

<br>

#### Check the tests for complete examples

https://github.com/hofstadter-io/hof/tree/_dev/test/render

<br>

#### Want to use and compose code gen modules and dependencies?

Create and use generator modules.

`hof gen app.cue -G frontend -G backend -G migrations`

See the [first-example](/first-example/) to learn how.

<br>

#### Command Help

<br>

{{<codePane file="code/cmd-help/gen" title="$ hof help gen" lang="text">}}

{{< childpages >}}

