---
title: "Command Help"
description: "Help text for hof's main commands"
brief: "Help text for hof's main commands"

weight: 100
---

`hof` has several commands typically run during daily development.
The `the-walkthrough` will introduce you to these commands
and each has a dedicated section for the topic as well.

### hof

Top level commands and help message

<details>
<summary>hof help</summary>
{{<codePane file="code/cmd-help/hof" title="$ hof help" lang="text">}}
</details>


### hof / gen

Declarative code generation for directories and files from data, CUE, and templates.
Build with adhoc one-liners or use composable generators to create reusable blueprints.

See [the code generation section](/code-generation/) to learn more

<details>
<summary>hof help flow</summary>
{{<codePane file="code/cmd-help/gen" title="$ hof help gen" lang="text">}}
</details>


### hof / create

Generate boilerplate from any git repository using hof.

See [the creators section](/code-generation/creators/) to learn more

<details>
<summary>hof help flow</summary>
{{<codePane file="code/cmd-help/gen" title="$ hof help gen" lang="text">}}
</details>


### hof / flow

Build workflows and scripts with CUE and a DAG engine

See [the task engine section](/task-engine/) to learn more

<details>
<summary>hof help flow</summary>
{{<codePane file="code/cmd-help/flow" title="$ hof help flow" lang="text">}}
</details>


### hof / datamodel

Used for data model management (dm for short)

See the [data modeling section](/data-modeling/) for details.

<details>
<summary>hof help datamodel</summary>
{{<codePane file="code/cmd-help/dm" title="$ hof help datamodel" lang="text">}}
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


### hof / fmt

Format many languages at once with good defaults.

See the [formatting](/code-generation/formatting) for details.

<details>
<summary>hof help fmt</summary>
{{<codePane file="code/cmd-help/fmt" title="$ hof help fmt" lang="text">}}
</details>


### hof / eval

CUE eval embedded in hof.

<details>
<summary>hof help eval</summary>
{{<codePane file="code/cmd-help/eval" title="$ hof help eval" lang="text">}}
</details>


### hof / export

CUE export embedded in hof.

<details>
<summary>hof help export</summary>
{{<codePane file="code/cmd-help/export" title="$ hof help export" lang="text">}}
</details>


### hof / vet

CUE vet embedded in hof.

<details>
<summary>hof help vet</summary>
{{<codePane file="code/cmd-help/vet" title="$ hof help vet" lang="text">}}
</details>
