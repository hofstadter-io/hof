---
title: "hof / flow"
linkTitle: "hof / flow"
weight: 35
---

{{<lead>}}
`hof/flow` is a data and task engine
with automatic dependency detection
powered by `cue/flow`.
It has more task types and capabilities.
{{</lead>}}

{{<codeInner title="hof flow -h" >}}
run file(s) through the hof/flow DAG engine

Use hof/flow to transform data, call APIs, work with DBs,
read and write files, call any program, handle events,
and much more.

'hof flow' is very similar to 'cue cmd' and built on the same flow engine.
Tasks and dependencies are inferred.
Hof flow has a slightly different interface and more task types.

Docs: https://docs.hofstadter.io/data-flow

Example:

  @flow()

  call: {
    @task(api.Call)
    req: { ... }
    resp: {
      statusCode: 200
      body: string
    }
  }

  print: {
    @task(os.Stdout)
    test: call.resp
  }

Arguments:
  cue entrypoints are the same as the cue cli
  @path/name  is shorthand for -f / --flow should match the @flow(path/name)
  +key=value  is shorthand for -t / --tags and are the same as CUE injection tags

  arguments can be in any order and mixed

@flow() indicates a flow entrypoint
  you can have many in a file or nested values
  you can run one or many with the -f flag

@task() represents a unit of work in the flow dag
  intertask dependencies are autodetected and run appropriately
  hof/flow provides many built in task types
  you can reuse, combine, and share as CUE modules, packages, and values

Usage:
  hof flow [cue files...] [@flow/name...] [+key=value] [flags]

Aliases:
  flow, f

Flags:
  -d, --docs           print pipeline docs
  -f, --flow strings   flow labels to match and run
  -h, --help           help for flow
  -l, --list           list available pipelines
      --progress       print task progress as it happens
  -s, --stats          Print final task statistics
  -t, --tags strings   data tags to inject before run

Global Flags:
  -p, --package string   the Cue package context to use during execution
  -q, --quiet            turn off output and assume defaults at prompts
  -v, --verbose int      set the verbosity of output
{{< /codeInner >}}


## args & flags

`hof/flow` accepts CUE entrypoints like the other commands.
There is CLI sugar for

- flows: `@path/name` is sugar for `-f path/name`
- tags:  `+key=value` is sugar for `-t key=value`

Flags:

- `-f`/`@` is used to select a flow by name in `@flow(<name>)`
- `-t`/`+` is used to inject strings into tags `@tag(<name>)`
- `-l`/`--list` prints the list of discovered flows
- `-d`/`--docs` prints additional flow details and docs
- `--progress` will print task progress for the events found, pre, & post
- `--stats` will print task times and dependencies at completion 


---

{{< childpages >}}
