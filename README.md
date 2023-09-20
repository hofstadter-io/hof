# hof - the higher order framework

A tool that unifies data models, schemas, code generation, and a task engine.

__`hof` is a CLI tool you *add* to your workflow.__

- Augment data, config, and schemas with CUE to improve consistency, gain confidence, and unlock new possibilities
- Generate multiple Yaml and JSON files in one-shot, from a CUE based source of truth
- Generate much of the application code, work directly in the output, regenerate without loss
- Explore data or config with the TUI, work with it using CUE in live-reload mode
- Run workflows with automatic task dependency inference, the right tasks are run in the right order


| Core Features | |
|:---                    |:-- |
| __code generation__    | Data + templates = _ (anything), technology agnostic |
| __data model__         | Define and manage data models, life-cycle features |
| __task engine__        | Extensible task and DAG workflow engine |
| __CUE cmds__           | Core def, eval, export, and vet commands |
| __creators__           | bootstraping and starter kits from any repo |
| __modules__            | CUE module dependency management |
| __tui__                | A terminal interface to Hof and CUE |
| __chat__               | Combine LLM and Hof code gen for better, scalable results |

<br>

`hof` uses [CUE](https://cuelang.org) to power the DX and implementation.
We believe CUE is a great language for specifying schemas, configuration, and generally
for writing anything declarative or this is a source of truth.
It has good theory and comes from the same people that brought us containers, Go, and Kubernetes.

<!-- something about osurce of thuth, unified abstraction later, interoperablility... -->

Learn more about CUE: [cuelang.org](https://cuelang.org) | [cuetorials.com](https://cuetorials.com)


## Documentation

Please see __[docs.hofstadter.io](https://docs.hofstadter.io)__ to learn more.

The [getting-started](https://docs.hofstadter.io/getting-started/) section will take you on a tour of hof.
The [first-example](https://docs.hofstadter.io/first-example/) section shows you how to build and use a generator.

Join us or ask questions on

- Slack: [https://hofstadter-io.slack.com](https://join.slack.com/t/hofstadter-io/shared_invite/zt-e5f90lmq-u695eJur0zE~AG~njNlT1A)
- Discord: https://discord.com/invite/BXwX7n6B8w

We also use GitHub issues and discussions. Use which every is easiest for you!


## Installation

You can find [the latest downloads on our GitHub releases page](https://github.com/hofstadter-io/hof/releases).
This is the preferred method.

If you already have hof, install a specific version with `hof update --version vX.Y.Z`.

```shell
# Homebrew
brew install hofstadter-io/tap/hof

# Shell Completions (bash, zsh, fish, power-shell)
echo ". <(hof completion bash)" >> $HOME/.profile
source $HOME/.profile

# Show the help text or version info to verify installation
hof --help
hof version
```


## Contributing & Community

Interested in helping out or hanging out? The best ways to get started are

1. [Joining the Community](https://docs.hofstadter.io/help-and-support/#community/)
1. [The Contributing Guild](https://docs.hofstadter.io/contributing/)


## Interfaces 

There are two interfaces to `hof`

1. a CLI - great for scripting and automation
2. a TUI - great for exploring and designing

A desktop version is in the works, reach out if you would like early access.

### cli

```
hof - the higher order framework

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]

Main commands:
  chat                  co-create with AI (alpha)
  create                starter kits or blueprints from any git repo
  datamodel             manage, diff, and migrate your data models
  def                   print consolidated CUE definitions
  eval                  evaluate and print CUE configuration
  export                output data in a standard format
  flow                  run workflows and tasks powered by CUE
  fmt                   format any code and manage the formatters
  gen                   CUE powered code generation
  mod                   CUE module dependency management
  tui                   a terminal interface to Hof and CUE
  vet                   validate data with CUE

Additional commands:
  help                  help about any command
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal
  feedback              open an issue or discussion on GitHub

Flags:
  -E, --all-errors           print all available errors
  -h, --help                 help for hof
  -i, --ignore-errors        turn off output and assume defaults at prompts
  -D, --include-data         auto include all data files found with cue files
  -V, --inject-env           inject all ENV VARs as default tag vars
  -I, --input stringArray    extra data to unify into the root value
  -p, --package string       the Cue package context to use during execution
  -l, --path stringArray     CUE expression for single path component when placing data files
  -q, --quiet                turn off output and assume defaults at prompts
  -d, --schema stringArray   expression to select schema to apply to data files
      --stats                print generator statistics
  -0, --stdin-empty          A flag that ensure stdin is zero and does not block
  -t, --tags stringArray     @tags() to be injected into CUE code
  -v, --verbosity int        set the verbosity of output
      --with-context         add extra context for data files, usable in the -l/path flag

Use "hof [command] --help / -h" for more information about a command.
```

### tui

The `hof tui` is a terminal based interface to Hof's features.
It has a built in help system and documentation.
The following YouTube video provides a tour.


[![Tour Hof's TUI](http://img.youtube.com/vi/XNBqBWO4y08/0.jpg)](http://www.youtube.com/watch?v=XNBqBWO4y08 "Hof TUI Overview")

