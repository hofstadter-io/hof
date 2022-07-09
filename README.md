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
  gen                   create arbitrary files from data with templates and generators
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

