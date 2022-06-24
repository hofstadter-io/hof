# hof - the high code framework

`hof` tries to remove redundent and error prone development activities.
The idea is to write your data model once and generate most of your code.
You can then work directly in the output to customize,
update the design or generators, and then regenerate your application.

`hof` uses CUE extensively to power the DX and implementation.

- data model management so you can checkpoint, diff, and migrate
- code generation to skaffold code across the stack
- diff3 for custom code and skaffold regeneration
- modular and composable code generators with dependency management

There are also several utilities subcommands for
adhoc file generation, data transformations, and
CUE powered scripting.

## Install

You will have to download `hof` the first time.
After that `hof` will prompt you to update and
install new releases as they become available.

You can find the latest version from the
[releases page](https://github.com/hofstadter-io/hof/releases)
or use `hof` to install a specific version of itself with `hof update --version vX.Y.Z`.

```shell
# Latest Release
go install github.com/hofstadter-io/hof/cmd/hof@latest

# Latest Commit
go install github.com/hofstadter-io/hof/cmd/hof@_dev

# Shell Completions (bash, zsh, fish, power-shell)
echo ". <(hof completion bash)" >> $HOME/.profile
source $HOME/.profile

# Show the help text
hof --help
```

## Documentation

Please see __https://docs.hofstadter.io__ to learn more.

The [first-example](https://docs.hofstadter.io/first-example)
will take you through the process
of creating and using a simple generator

Join us on Slack! [https://hofstadter-io.slack.com](https://join.slack.com/t/hofstadter-io/shared_invite/zt-e5f90lmq-u695eJur0zE~AG~njNlT1A)

## Main Commands

### hof

```
hof - the high code framework

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]

Main commands:
  datamodel             manage, diff, and migrate your data models
  render                generate arbitrary files from data and CUE entrypoints
  gen                   render directories of code using modular generators
  flow                  run CUE pipelines with the hof/flow DAG engine
  mod                   polyglot dependency management based on go mods and MVS

Additional commands:
  help                  help about any command
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal
  feedback              send feedback, bug reports, or any message

Flags:
  -h, --help             help for hof
  -p, --package string   the Cue package context to use during execution
  -q, --quiet            turn off output and assume defaults at prompts
  -v, --verbose int      set the verbosity of output

Use "hof [command] --help / -h" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.
```


### datamodel

```
Data models are sets of models which are used in many hof processes and modules.

At their core, they represent the most abstract representation for objects and
their relations in your applications. They are extended and annotated to add
context fot their usage in different code generators: (DB vs Server vs Client).

Beyond representing models in their current form, a history is maintained so that:
  - database migrations can be created and managed
  - servers can handle multiple model versions
  - clients can implement feature flags
Much of this is actually handled by code generators and must be implemented there.
Hof handles the core data model definitions, history, and snapshot creation.

Usage:
  hof datamodel [command]

Aliases:
  datamodel, dm

Available Commands:
  checkpoint  create a snapshot of the data model
  diff        show the diff between data model version
  history     list the snapshots for a data model
  info        print details for a data model
  list        print available data models
  log         show the history of diffs for a data model

Flags:
  -d, --datamodel strings   Datamodels for the datamodel commands
  -f, --format string       Pick format from Cuetils (default "_")
  -h, --help                help for datamodel
  -m, --model strings       Models for the datamodel commands
  -o, --output string       Output format [table,cue] (default "table")
  -s, --since string        Timestamp to filter since
  -u, --until string        Timestamp to filter until

Global Flags:
  -p, --package string   the Cue package context to use during execution
  -q, --quiet            turn off output and assume defaults at prompts
  -v, --verbose int      set the verbosity of output

Use "hof datamodel [command] --help" for more information about a command.
```

### render

```
hof render joins CUE with Go's text/template system and diff3
  create on-liners to generate any file from any data
  edit and regenerate those files while keeping changes

# Render a template
hof render data.cue -T template.txt
hof render data.yaml schema.cue -T template.txt > output.txt

# Add partials to the template context
hof render data.cue -T template.txt -P partial.txt

# The template flag
hof render data.cue ...

  # Multiple templates
  -T templateA.txt -T templateB.txt

  # Cuepath to select sub-input values
  -T 'templateA.txt;foo'
  -T 'templateB.txt;sub.val'

  # Writing to file
  -T 'templateA.txt;;a.txt'
  -T 'templateB.txt;sub.val;b.txt'

  # Templated output path 
  -T 'templateA.txt;;{{ .name | ToLower }}.txt'

  # Repeated templates when input is a list
  #   The template will be processed per item
  #   This also requires using a templated outpath
  -T 'template.txt;items;out/{{ .filepath }}.txt'

# Learn about writing templates, with extra functions and helpers
  https://docs.hofstadter.io/code-generation/template-writing/

# Check the tests for complete examples
  https://github.com/hofstadter-io/hof/tree/_dev/test/render

# Want to use and compose code gen modules and dependencies?
  hof gen is a scaled out version of this command
  hof gen app.cue -g frontend -g backend -g migrations

Usage:
  hof render [flags] [entrypoints...]

Aliases:
  render, R

Flags:
  -D, --diff3              enable diff3 support, requires the .hof shadow directory
  -h, --help               help for render
  -P, --partial strings    file globs to partial templates to register with the templates
  -T, --template strings   Template mappings to render with data from entrypoint as: <filepath>;<?cuepath>;<?outpath>

Global Flags:
  -p, --package string   the Cue package context to use during execution
  -q, --quiet            turn off output and assume defaults at prompts
  -v, --verbose int      set the verbosity of output
```

### gen

```
render directories of code using modular generators

  https://docs.hofstadter.io/first-example/

hof gen app.cue -g frontend -g backend -g migrations

Usage:
  hof gen [files...] [flags]

Aliases:
  gen, G

Flags:
  -g, --generator strings   Generators to run, default is all discovered
  -h, --help                help for gen
  -s, --stats               Print generator statistics

Global Flags:
  -p, --package string   the Cue package context to use during execution
  -q, --quiet            turn off output and assume defaults at prompts
  -v, --verbose int      set the verbosity of output
```

### flow

```
run CUE pipelines with the hof/flow DAG engine

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
```
