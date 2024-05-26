---
title: Task Engine
weight: 50
---

{{<lead>}}
Hof's Task Engine is a CUE powered DAG processor.
Define tasks and flows in CUE to process data,
implement webhooks, or connect APIs.
Use the library of builtin tasks and connectors
or create custom tasks using any container.
Hof will infer dependencies and run tasks
as they are ready or needed.
{{</lead>}}


{{< childpages >}}


### Attributes

`@flow()`: useful for a single, unnamed workflow

`@flow(<name>)`: write and call one or more flows by name

`@task(<pkg.Func>)`: a workflow step of type `<pkg.Func>`

`@print(<cue.path>)`: output values to log for debugging

### Command Help

<details>
<summary>hof flow -h</summary>
{{<codePane title="hof flow -h" file="code/cmd-help/flow" lang="text">}}
</details>


### args & flags

`hof/flow` accepts CUE entrypoints like the other commands.
There is CLI sugar for

- flows: `@path/name` is sugar for `-F path/name`
- tags:  `+key=value` is sugar for `-t key=value`

Useful Flags:

- `-F`/`@` is used to select a flow by name in `@flow(<name>)`
- `-t`/`+` is used to inject strings into tags `@tag(<name>)`
- `--progress` will print task progress for the events found, pre, & post

