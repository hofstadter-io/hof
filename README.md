# hof - the high code framework

`hof` combines data models, code generation, and modules
to help you write and maintain large amounts of code.

1. __data model__ - define & manage data models - the source of truth
2. __code generation__ - data + template = _ (anything) - technology agnostic
3. __modules__ - composable data models & generators - an ecosystem

<img src="./images/how-hof-works.svg" alt="how hof works" width="100%" height="auto" style="max-width:600px">

__`hof` is a CLI tool you will add to your workflows.__
It is technology agnostic, captures common patterns and boilerplate,
has modules that span technologies, and continues to work as your application evolves.

- data model management so you can checkpoint, diff, and calculate migrations
- code generation to scaffold consistent code and boilerplate across the stack
- diff3 to support custom code, data model updates, and code regeneration
- modular and composable code generators with dependency management

`hof` uses [CUE](https://cuelang.org) extensively to power the DX and implementation.
There are also several utilities subcommands for
adhoc file generation, data transformations, and
CUE powered scripting.

## [try hof on github](https://github.com/codespaces/new?hide_repo_select=true&ref=main&repo=604970115)

## Install

You will have to download `hof` the first time.
After that `hof` will prompt you to update and
install new releases as they become available.

You can find the latest version from the
[releases page](https://github.com/hofstadter-io/hof/releases)
or use `hof` to install a specific version of itself with `hof update --version vX.Y.Z`.

```shell
# Homebrew
brew install hofstadter-io/tap/hof

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

Please see __[docs.hofstadter.io](https://docs.hofstadter.io)__ to learn more.

The [first-example](https://docs.hofstadter.io/first-example)
will take you through the process
of creating and using a simple generator
[Several demos](https://github.com/hofstadter-io/demos) in a separate repository
and various `hofmod-*` repositories are available for you to use.

Join us on Slack! [https://hofstadter-io.slack.com](https://join.slack.com/t/hofstadter-io/shared_invite/zt-e5f90lmq-u695eJur0zE~AG~njNlT1A)
We are more than happy to answer your questions.


## Main Commands

### hof

```
hof - the high code framework

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]

Main commands:
  create                bootstrap projects, components, and files from any git repo
  datamodel             manage, diff, and migrate your data models
  gen                   modular and composable code gen: CUE & data + templates = _
  flow                  run CUE pipelines with the hof/flow DAG engine
  fmt                   format any code and manage the formatters
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
  -v, --verbosity int    set the verbosity of output

Use "hof [command] --help / -h" for more information about a command.
```

