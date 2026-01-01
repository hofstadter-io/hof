# hof/env

`hof/env` is a CUE interface over Dagger.
The result is a blend of Makefiles, Ansible, Dockerfiles, and Docker Compose.

It forms the foundation for:

1. Shared fabrics across the software lifecycle, via CUE and Dagger DAGs
1. OCI modules and imports in the front, OCI images and environments in the back
1. Defining images/layers, services/deployments, and workflows/tasks
1. Works the same everywhere with dynamic tasks dependencies and advanced caching

> [!NOTE]
> A name change is approaching, `hof` -> `veg`, you will see hints of that now.

## Getting Started

You need docker, nerdctl, or podman installed.

<!-- ### Install with Homebrew

```sh
brew install hofstadter-io/tap/hof
``` -->

```sh
hof env install/hof

dagger ...

docker copy ghcr.io

curl github
```

### Binaries from GitHub

|                                              eat your                                               |                                              veggies                                               |
| :-------------------------------------------------------------------------------------------------: | :------------------------------------------------------------------------------------------------: |
| [linux / arm](https://github.com/hofstadter-io/hof/releases/download/v0.7.0/hof_v0.7.0_linux_arm64) | [mac / arm](https://github.com/hofstadter-io/hof/releases/download/v0.7.0/hof_v0.7.0_darwin_arm64) |
| [linux / amd](https://github.com/hofstadter-io/hof/releases/download/v0.7.0/hof_v0.7.0_linux_amd64) | [mac / amd](https://github.com/hofstadter-io/hof/releases/download/v0.7.0/hof_v0.7.0_darwin_amd64) |

### Run an example

```sh
# get the repo
git clone https://github.com/hofstadte-io/hof && cd hof
make registry.start # todo, hof should handle this like it does for formatters
make dagger.start

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

> [!INFO]
> `schemas/env` and `lib/env/dag` align very closely. The `dag` package uses a new progressive CUE schema alignment and decoding strategy
> that works really, really well and will be used more widely in veg. It's also prime time for `veg gen`.

### The `veg/env` command

The `veg/env` command aims to be flexible, extensible, and consistent

- All take the same flags and select / filter the same way

Format: `veg env [cmd] [flags] [args] % [cue entrypoints]`

- `cmd` is a built-in or custom command to run. Use `cmd/task` to subselect. Both are regexp.
- `flags` there are two group,
  - `-K/--kind` and `-P/--path` combine with `[args] to select targets. (todo, incorp label system)
  - `--no-cache`, `--no-exit`, `
- `args` are a regexp match on names
  - each arg is processed sequentially and independently
  - any args after a `%` are considered entrypoint paths to the underlying CUE evaluator

Guidance on starting out with the commands

- `veg/env` tries to be context aware
  - knows the project or subdirectory you are in
  - knows the object type when handling args (veg-dist as example)
- the commands use the same flags and filtering
  - start with list, then run and a workful command
  - start with build or export before publish
- use `env.Terminal` or `-F/-N` to go interactive
- import `schemas/env/rrr:env` instead of `schemas/env` for stricter schemas. It's slower because of the recursive disjunctions, but it will catch errors earlier and narrow down messages, especially with the aliasv2 experiment we are relying on.

> [!NOTE] > `veg env` without any subcommands will run your custom commands

#### Help Text

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

Keep [schema/env](../../schemas/veg)
and [lib/env/common](../../lib/env/common)
handy for the details of the following.

Generally speaking...

- There are several groups or classes of statements, all prefixed by the `schemas/env.*` package identifier.
  - `env.Step` maps onto `With<Step>` and `Without<Step>` and can appear in `#Container: steps: [...]`
  - `env.#Stuff` maps onto resources, artifacts, and Dagger types. They are inputs, intermediates, outputs, or runnable.
  - `env.$Func` maps from one resource to one of the same or another, some `env.#Stuff` do some of these naturally too.
- It's a one way trip from CUE -> Dagger, you cannot for instance, use a directory listing or http response in CUE
  - `hof/flow` exists for this use case and some merging is on the roadmap.
  - The key requirement to maintain distinct operation modes, hermetic and yolo, with control over where, when, and how.

> [!WARNING]
> We have swapped the semantics to `Token` | `#Token` from `WithToken()` and `Token()` with Dagger.
>
> 1. veg: `Dir` ~ dag: `WithDir()`
> 2. veg: `#Dir` ~ dag: `Dir()`

### Steps:

```
Dir               add a directory
File              add a file
EnvVar            add a single env var
SecretVar         add a single secret
EnvFile           add an env var file
SecretFile        add a secret var file

Temp              a temp volume for the next exec
Mount             mount a cache, file, directory, secret
BindService       bind another service to the container  (hint, dep graph)
Expose            mark a port for servin

Sync              force evaluation of the dagger graph
Exec              run any command as a container layer
Sh, Bash, Zsh     exec wrappers with a 'script' param
User              set the current user
Workdir           set the current workdir


Entrypoint        set container entrypoint
DefaultArgs       set container default args
DefaultTerm       set the terminal dagger uses when needed

!!!
Terminal          drop into a terminal at any or many point(s), directory or container
!!!               (this is one of the coolest dagger features)

More to come...

- VsCode (like terminal, combo of them too)
- Chown
- $Filter (#Dir->#Dir)
- $Diff (#Dir-#Dir->#Dir)
- Patch & #Patch
- Diff & #Diff, Changes & #Changeset
- ?Merge (not overwrite, doesn't exist yet)
- #DirToGit         git from a dir

Vscode            (we'll add a Step to open in vscode, or make something that does both, configurablely)

Without...
  -- both #Dir and #Container
  Dir
  File
  Files
  -- #Containers only
  EnvVar
  SecretVar

```

### #Stuff:

These are artifacts, intermediates, or resources you can work with

```
#DockerBuild      what you would expect
#Container        this is the (dagger) way
#Service          configure a container for running exposed

#Dir              a dir that can be used in CUE
#File             a file that can be used in CUE
#Secret           a secret that will be elided from output
#Cache            a named volume in memory, persists sessions

#Cmd              a custom command with named subtasks
#Task             a task is a list of runnables and is runnable itself

#GitRepo          from a uri
#HostDir          from a path
#HostFile         from a path
#HostImageFile    from a tarball
#HostImage        from local engine
#HostService      expose host to dagger
#HostTunnel       expose dagger to host
#HostSocket       from a path

#ExportDir        to host path
#ExportFile       to host path
#ExportImageFile  to a tarball
#ExportImage      to local engine
#PublishImage     to a registry

# many data formats available
#ExportCuefig     returns a #File for the CUE representation any #Thing
#ExportDagger     returns a #File for the Dagger representation any #Thing
```

## Examples

### Build and Run a Container

`hof env ...`

```cue
dev: env.#Container & {
  @env(), @id(veg-dev)
  #hof: metadata: {
    id:          "veg-dev"
    name:        id
    description: "setup needed to work on veg"
  }
  name: #hof.metadata.name

  // start from our default debian
  from: bases.debian.default

  steps: [
    // customization
    tool.zsh.customize,

    // daily drivers
    hof.cli,
    tool.github.cli,

    // deps for go/node/python -> c/c++ situations (like CGO)
    utils.apt.install & {#pkgs: ["gcc", "libc6-dev"]},

    // setup languages
    lang.go.default,
    lang.cue.default,
    lang.node.default,
    lang.python.default,
    lang.python.dev, // depends on node

    // devops stuff
    tool.hashicorp.terraform,
    tool.hashicorp.packer,
    tool.k8s.kubectl,
    tool.k8s.helm,
    tool.k8s.crane,

    // tools just for agents
    tool.agents.lsp2mcp,

    // launch a terminal to check on things, continue when done
    // env.Terminal,

    // bind lsp servers, started on demand
    env.BindService & {service: lang.go.lsp},
    env.BindService & {service: lang.cue.lsp},
    env.BindService & {service: lang.node.lsp},
    env.BindService & {service: lang.python.lsp},

  ]
}
```

### Multi-Stage Builds and Beyond, DAG style

- multi-stage
- no need for yum rm
- how binaries and dirs work

That long-time advice to install packages like this: `apt update && apt install && apt clean`... it's over!
We can now attach caches, just like we do for languages like `go.mod` and `node_modules`,
to save context and time while keeping images clean and slim.

`hof env ...`

```cue
apt: {
	caches: {
		varLib: env.#Cache & {
			name: "debian-13-var-lib-cache"
		}
	}

	mounts: {
		varLib: env.Mount & {
			path:   "/var/lib/apt/lists"
			source: apt.caches.varLib
		}
	}

	// generalized apt package install
	install: env.Bash & {
		#pkgs: [...string]
		script: "apt-get install -y --no-install-recommends \(strings.Join(#pkgs, " "))"
	}

	// runs apt-get update, do this once early
	update: env.Bash & {script: "apt-get update -y"}

	// You should NEVER need this again!
	// we use caches to do even better than either method
	// 1. same size savings as ( [update -> install -> clean] )
	// 2. save time with cache ( update -> [install] ... magic)
	// anyway, it cleans apt stuff
	clean: env.Bash & {
		script: """
			apt-get dist-clean
			rm -rf /var/lib/apt/lists/*
			"""
	}
}
```

You can then create a minimal base image with this mounted.
Every step afterwards will use this mount while in Dagger
and exported when exported they are not included.

`hof env ...`

```cue
	minimal: env.#Container & {
		@id(debian-13-minimal)
		#hof: metadata: description: "A minimal debian13 image with updates and certs"

		from: "debian:13-slim"

		steps: [
			// default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},

      // staying clean caches
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "debian-13-var-log"}},
			env.Mount & {path: "/var/cache", source: env.#Cache & {name: "debian-13-var-cache"}},

			// Shared, persistent Apt caches, for all derived images too
			// ...instead of cleaning and refetching all the time? (we like pain in devops #yamhell)
			utils.apt.mounts.varLib,

			// Update once at the beginning
			utils.apt.update,

			// just certs
			utils.apt.install & {#pkgs: ["ca-certificates"]}, // shouldn't need wget/curl, we can do that at this level
		]
	}
```

### Building with an existing Dockerfile

You can still build using your existing Dockerfiles.
It's also easy to patch source code or use the resulting image anywhere `hov/env`.

`hof env ...`

```cue
// All your code belong to us
repos: {
  blebbit: env.#GitRepo & {url: "https://github.com/blebbit/atproto"}
  atproto: env.#GitRepo & {url: "https://github.com/bluesky-social/atproto"}
  didplc: env.#GitRepo & {url: "https://github.com/did-method-plc/did-method-plc"}
  indigo: env.#GitRepo & {url: "https://github.com/bluesky-social/indigo"}
  jetstream: env.#GitRepo & {url: "https://github.com/bluesky-social/jetstream"}
}
// Permissioned PDS
ppds: {
  code: env.#Dir & {sources: [repos.blebbit]}
  ctr: env.#DockerBuild & {source: code, dockerfile: "services/pds/Dockerfile"}
}
// Official PDS
pds: {
  code: env.#Dir & {sources: [repos.atproto]}
  ctr: env.#DockerBuild & {source: code, dockerfile: "services/pds/Dockerfile"}
}
// Official PLC (with a patch for CNPG friendly db conn strings)
plc: {
  code: env.#Dir & {sources: [repos.didplc]}
  // patch
  fixd: env.#Dir & {sources: [repos.didplc], patch: patches.plc}
  ctr: env.#DockerBuild & {source: fixd, dockerfile: "packages/server/Dockerfile"}
}
// Official Relay (with a patch to remove git info in go build)
relay: {
  code: env.#Dir & {sources: [repos.indigo]}
  // patch
  fixd: env.#Dir & {sources: [repos.indigo], patch: patches.relay}
  ctr: env.#DockerBuild & {source: fixd, dockerfile: "cmd/relay/Dockerfile"}
}
// Official Jetstream
jetstream: {
  code: env.#Dir & {sources: [repos.jetstream]}
  ctr: env.#DockerBuild & {source: code}
}
```

### Commands in an Environment

`hof env ...`

```cue
cmd: {
	// unify in `@env()` and `name` two levels deep
	// with [patternMatching]~(keyAlias,_valAlias): { name: keyAlias }
	[string]~(k1,_): env.#Cmd & {
		@env(), name: k1
		tasks: [string]~(k2,_): {
			@env(), name: k2
			steps: [...[...{name: "\(k1).\(k2)"}]]
		}
	}

	test: tasks: {
		go: steps: [[_tester & {#cmd: "go test ./..."}]]
		// parallel tests
		goUltra: steps: [[
			_tester & {#cmd: "go vet ./..."},
			_tester & {#cmd: "go test -race ./..."},
			_tester & {#cmd: "go test -cover ./..."},
		]]
		// sequential tests
		// vet: {steps: [[_tester & {#cmd: "go vet ./..."}]]}
		// race: {steps: [[_tester & {#cmd: "go test -race ./..."}]]}
		// cover: {steps: [[_tester & {#cmd: "go test -cover ./..."}]]}
	}
	lint: tasks: {
		// want something like: gofmt -l . | wc -l | grep -e '^0$'
		fmt: steps: [[_tester & {#cmd: #"gofmt -l . || true"#}]]
		staticcheck: steps: [[_tester & {#cmd: "staticcheck ./... || true"}]]
		golangci: steps: [[_tester & {#cmd: "golangci-lint run || true"}]]
		spelling: _
	}

	scan: tasks: {
		sonar: {}
		vuln: {}
	}

	review: tasks: {
		agent: {
			... code changes,
			docs / agents.md need updating,
			stage & apply suggested changes,
		}
	}

	ci: tasks: {
		default: steps: [test, lint, scan]
		onPush:  steps: [test, lint]
		prPush:  default

		prepare: [...]
		release: [ci.default, prepare]
		onTag:   [release]
	}
}

```

### Release Bundles

You can define release bundles and then assemble and publish them with a single command.

`hof env ...`

```

```

### Agent or Dev Environments with Tools

[veg-dev]

(test,lint,lsp) - stuff to make things easier, show service bindings

### Webhooks and Server Mode

## Patterns

todo:

- working with files and directories
- the art of "container" composition (beyond multi-stage)
- insert `env.Terminal` to debug, use flags too

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

### bring your dotfiles and customization to any image

### Flags, Configs, and Defaults

- using flags to override
- using data to override
- parameterize large swaths, show the propagation

### Matrix Comprehension

CUE has list and struct comprehension we can use to

- github like matrix [os,lang,version] CI
- multi-arch images and binaries
- container families and parameterization

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
