---
title: Task Engine
description: "Build and run CUE based task workflows"
brief: "Build and run CUE based task workflows"

weight: 35
---

{{<lead>}}
Hof's __task engine__ allows you to write dynamic workflows with CUE.
Build modular and composable DAGs that can

- work with files, call APIs, can prompt for user input, and much more
- run scripts and containers, use any language for a custom task
- run a server that calls a workflow per request or bulk process a set of inputs
- create additional task based on the results of other tasks

Hof's task engine is an extension of cue/flow with

- extra tasks for key-value cache and advanced structural operations on CUE values
- additional control over parallel execution, synchronization of tasks, and message passing
- allows for workflows along side regular CUE files so they can be imported and shared as modules
{{</lead>}}


## Overview


### Command

{{<codeInner title="hof flow help">}}
run workflows and tasks powered by CUE

Usage:
  hof flow [cue files...] [@flow/name...] [+key=value] [flags]
  hof flow [command]

Aliases:
  flow, f

Available Commands:
  list        print available flows

Flags:
  -B, --bulk string        exprs for inputs to run workflow in bulk mode
  -F, --flow stringArray   flow labels to match and run
  -h, --help               help for flow
  -P, --parallel int       global flow parallelism (default 1)
      --progress           print task progress as it happens

Global Flags:
  ...
{{</codeInner>}}

### Tasks & Schemas

You can find the schema and example for all tasks in
[the hof/flow reference section](/reference/hof-flow/schemas/)

- `api`
  - Call
  - Serve
- `csp` (communicating sequential processes)
  - Chan
  - Send
  - Recv
- `cue`
  - Format (print incomplete to concrete CUE values)
- `gen` (generate random values)
  - Seed
  - Now
  - Str
  - Int
  - Float
  - Norm
  - UUID
  - CUID
  - Slug
- `hof`
  - Template (render a hof text/template)
- `kv`
  - Mem (in-memory cache)
- `os`
  - Exec
  - FileLock
  - FileUnlock
  - Getenv
  - Glob
  - Mkdir
  - ReadFile
  - ReadGlobs
  - Sleep
  - Stdin
  - Stdout
  - Watch
  - WriteFile
- `prompt`
  - Prompt (interactive user prompts, like creators)
- `st` (structural)
  - Mask
  - Pick
  - Insert
  - Replace
  - Upsert
  - Diff
  - Patch

## Examples

You can find many examples through hof's codebase and other projects

- [the hof/flow reference section](/reference/hof-flow/schemas/)
- [the hof/flow test directory](https://github.com/hofstadter-io/hof/tree/_dev/flow/testdata)
- [LLM chat examples](https://github.com/hofstadter-io/hof/tree/_dev/flow/chat)
- [event and server based processing](https://github.com/verdverm/streamer-tools)
