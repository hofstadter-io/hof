# hof/env

`hof/env` is a CUE interface over Dagger.
The result is a blend of Makefiles, Dockerfiles, and Docker Compose.

It forms the foundation for:

1. Shared fabrics across the life-cycle, both CUE and Dagger DAGs
1. OCI modules and imports in the front, OCI images and environments in the back
1. Defining images/layers, services/deployments, and workflows/tasks
1. Works the same everywhere, the same `hof env ...` commands built with the same containerized workflows.

## Getting Started

<!-- ### Install with Homebrew

```sh
brew install hofstadter-io/tap/hof
``` -->

### Binaries from GitHub

|  eat your                                                                                           |  veggies                                                                                           |
|:---------------------------------------------------------------------------------------------------:|:--------------------------------------------------------------------------------------------------:|
| [linux / arm](https://github.com/hofstadter-io/hof/releases/download/v0.7.0/hof_v0.7.0_linux_arm64) | [mac / arm](https://github.com/hofstadter-io/hof/releases/download/v0.7.0/hof_v0.7.0_darwin_arm64) |
| [linux / amd](https://github.com/hofstadter-io/hof/releases/download/v0.7.0/hof_v0.7.0_linux_amd64) | [mac / amd](https://github.com/hofstadter-io/hof/releases/download/v0.7.0/hof_v0.7.0_darwin_amd64) |

### Run an example

```sh
# get the repo
git clone https://github.com/hofstadte-io/hof && cd hof

# pick an example
cd examples/env/...

# poke around
hof env list
hof env list ['^name$'] [-K '^kind$'] [-P '^path$'] [-S name|kind|path]

# you can just do things!
hof env [init, run, build, export, publish, up]

# make your own commands and flags
hof env [init, sync, test, lint, scan, promote, deploy, ci, ...]
hof env ... -t env=stg -t stack=app -t branch=main

# one-liners
cd veg       && hof env run   veg-dev
cd veg       && hof env build veg-ops
cd adk       && hof env test lint scan  # or `hof env ci`
cd atproto   && hof env up
cd inception && hof env turtles

# is k8s in dagger still a one-liner?
cd k8s && \
   hof env init && \
   hof env boot && \
   hof env lgtm && \
   hof env helm && \
   hof env test

```

### Code Organization

#### CUE (user) files

[todo] make these links

- `schemas/env/`
- `lib/env/common/` (these are import ordered to avoid cycles)
  - `utils/` helpers that only import veg/schemas/...
  - `bases/` operating systems and the like
  - `steps/` something like ansible / multi-stage dockerfile
  - `packs/` abstractions, collections, and such for reuse
- `examples/env/`
  - `basic/` 3-tier app
  - `adk/` matrix test & lint
  - `atproto/` compose like testnet & app
  - `veg/` monorepo devx & ci
  <!-- - `gitops/`      ordering tf & helm better
  - `inception/` various nesting setups -->

#### Go (impl) files:

- `lib/env/...`

### Help Text

`hof env` without any subcommands will run your commands

```
build, run, ship, and deploy environments (image, service, stack)

Usage:
  hof env [args] [flags]
  hof env [command]

Available Commands:
  build       build an environment
  export      export an environment into local container runtime
  get         get details for an environments
  list        list environments
  publish     publish an environment
  run         run an interactive environment
  up          starts an environment

Flags:
  -h, --help               help for env
  -K, --kind stringArray   kinds to include, defaults to all
  -Z, --no-cache           bust the cache and force evaluation
  -N, --no-exit            Leave the TUI open after finishing
  -F, --on-failure         on failure, enter an interactive terminal, requires a tty
  -P, --path stringArray   (cue) path prefixes to include, defaults to all
  -R, --renderer string    output format [auto, plain, tty, dots, report (for ai)] (default "auto")
  -S, --sort stringArray   sort columns, can be used multiple times (default [name])
```

## Steps and #Stuff

Keep `schema/env` and `lib/env/common` handy to see the details of each of these

- `env.Step` maps onto `With<Step>`
- `env.#Stuff` maps onto resources, artifacts, and Dagger types

### Steps:

```

```

### #Stuff:

```
#Container        can be run or published
#Service          can be up'd or attached

#File
#Dir
#Secret

#GitRepo          from a uri
#HostFile         from a path
#HostDir          from a path
#HostImage        from local engine
#HostImageFile    from a tarball

#ExportFile       to host path
#ExportDir        to host path
#ExportImage      to local engine
#ExportImageFile  to a tarball
#PublishImage     to a registry
```

## Examples

### Build and Run a Container

### Multi-Stage Builds, DAG style

### Building with an existing Dockerfile

[atproto]

(you can patch too)

### Commands in an Environment

[adk]

### Agent or Dev Environments with Tools

[veg-dev]

(test,lint,lsp) - stuff to make things easier, show service bindings

### Webhooks and Server Mode

### GitOps Images

Create a base image and family of specializations. Need to ship gitops containers to customers across clouds? Add an extra matrix dimension with another for loop.
`hof env build` will create them all without any arguments or use flags and args to parameterize, filter, and select a set of `#Things` to `hof env <op>`erate on or against.

```cue
// base gitops container
"ops": env.#Container & {
  from: bases.debian.minimal
  steps: [
    hof.cli,
    tool.hashicorp.terraform,
    tool.k8s.kubectl,
    tool.k8s.helm,
    tool.k8s.crane,
    tool.github.cli,
  ]
}

// variation per cloud to stay minimal
for c, cli in _clis {
  "ops-\(c)": env.#Container & {@env(), from: root.ctr.ops, steps: [cli]}
}

// bundled for multi or cross-cloud operations
"ops-all": env.#Container & {from: root.ctr["ops"], steps: [for _, cli in _clis {cli}]}

// this is a one-dimensional CI matrix
_clis: {
  gcp: tool.cloud.gcloud
  aws: tool.cloud.awscli
  az:  tool.cloud.azure
}
```

## Patterns

### custom steps with a single line

```sh
// https://learn.microsoft.com/en-us/cli/azure/install-azure-cli-linux?view=azure-cli-latest&pivots=apt
azureCli: env.Bash & { _script: "curl -sL https://aka.ms/InstallAzureCLIDeb | bash" }

image: {
  from: base.image
  steps: [
    azureCli,
    ...
  ]
}
```

### Flags, Configs, and Defaults

- using flags to override
- using data to override
- parameterize large swaths, show the propagation

## Notes

Use docker format & filters to inspect images

```sh
# view list & sort
docker image list --format "table {{.Repository}}:{{.Tag}}\t{{.Size}}" --filter "reference=veg-*" | grep -e '^veg' | sort -k2 -h

# remove matching patterns
docker rmi -f $(docker image list --format 'table {{.Repository}}:{{.Tag}}' | grep -e '^veg')

# inspect layers
dive veg-dev:local
```

What we use this CUE + Dagger magic for:

1. Powering a vscode virtual filesystem and diff viewer
1. Getting in on all that agentic hype, especiall tools, skills, and having...
   1. a safe space to work instead of restrictions being forced down
   1. git like ops on sessions and their environments
   1. fully recorded and sharable history via OCI
   1. Powering a Copilot alternative
   1. We now use this CUE + Dagger for both organic and agentic coding
1. Soon(?) CI, because there has to be a better way than Jenkins, Argo, GHA
1. There exist ambitions to tame terraform / helm sequencing and reconciliation
