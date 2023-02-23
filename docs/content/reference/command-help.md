---
title: "Command Help"
description: "Help text for hof's main commands"
brief: "of the hof CLI"

weight: 100
---

`hof` has several commands typically run during daily development.
The `first-example` will introduce you to these commands
and each has a dedicated section for the topic as well.

### hof

Top level commands and help message

<details>
<summary>hof help</summary>
{{<codePane file="code/cmd-help/hof" title="$ hof help" lang="text">}}
</details>


### hof / datamodel

Used for data model management (dm for short)

See the [Data Modeling section](/reference/hof-datamodel/) for details.

<details>
<summary>hof help datamodel</summary>
{{<codePane file="code/cmd-help/dm" title="$ hof help datamodel" lang="text">}}
</details>

### hof / gen

Create one-liners to generate files with data, CUE, and templates
or use composable generators to build out advanced applications.

See [the code gen docs](/reference/hof-gen/) to learn more

<details>
<summary>hof help flow</summary>
{{<codePane file="code/cmd-help/gen" title="$ hof help gen" lang="text">}}
</details>

### hof / flow

Build workflows and scripts with CUE and a DAG engine

See [the flow docs](/reference/hof-flow/) to learn more

<details>
<summary>hof help flow</summary>
{{<codePane file="code/cmd-help/flow" title="$ hof help flow" lang="text">}}
</details>


### hof / mod

Manage a MVS (Golang) style modules and dependencies.
Create custom module systems with a single config file.

Also, more typically, used as a stop gap for CUE modules until `cue mod` is implemented.

{{<codeInner title="typical usage">}}
# initialize a new module
hof mod init cue hof.io/docs/example

# download dependencies
hof mod vendor cue
{{</codeInner>}}

<details>
<summary>hof help mod</summary>
{{<codePane file="code/cmd-help/mod" title="$ hof help mod" lang="text">}}
</details>

