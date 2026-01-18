# veg/env

`veg/env` is a CUE interface over Dagger.
The result is a blend of Makefiles, Ansible, Dockerfiles, and Docker Compose.

This serves as a foundations for:

1. Shared configuration mesh across the software lifecycle, via CUE and Dagger DAGs
1. OCI modules and config in the front, OCI layers and environments in the back
1. Automatic DAG evaluation with dependency ordering and advanced caching

Define everything that goes into images/layers, stacks/services, and workflows/tasks which span local dev, ci/cd, and staged deployment. A singular, modular, shared fabric that acts as a source of truth and consistency from coder to commit to cloud.

> [!NOTE]
> A name change is approaching, `hof` -> `veg`, you will see hints of that now.

## Getting Started

You need docker, nerdctl, or podman installed.

<!-- ### Install with Homebrew

```sh
brew install hofstadter-io/tap/hof
``` -->

<!-- ```sh
hof env install/hof

dagger ...

docker copy ghcr.io

curl github
``` -->

### Binaries from GitHub

|                                              eat your                                               |                                              veggies                                               |
| :-------------------------------------------------------------------------------------------------: | :------------------------------------------------------------------------------------------------: |
| [linux / arm](https://github.com/hofstadter-io/hof/releases/download/v0.7.0-alpha.2/hof-linux-arm64) | [mac / arm](https://github.com/hofstadter-io/hof/releases/download/v0.7.0-alpha.2/hof-darwin-arm64) |
| [linux / amd](https://github.com/hofstadter-io/hof/releases/download/v0.7.0-alpha.2/hof-linux-amd64) | [mac / amd](https://github.com/hofstadter-io/hof/releases/download/v0.7.0-alpha.2/hof-darwin-amd64) |


### First Examples (low-level steps)

```cue
package basic

import "github.com/hofstadter-io/hof/schemas/env"

multi: {
	#ver:  string | *"1.25" @tag(ver)
	#arch: string           @tag(arch,var=arch)
	#os:   string           @tag(os,var=os)

	// base container for building in
	base: env.#Container & {
		@env()
		from: "golang:\(#ver)-alpine"
		steps: [
			// mount caches for mods and intermediate build artifacts (saves time)
			env.Mount & {path: "/cache/go", source: env.#Cache & {name: "go-build-\(#ver)-\(#arch)"}},
			env.Mount & {path: "/go", source: env.#Cache & {name: "go-mods-\(#ver)-\(#arch)"}},

			// set any default Go vars for all builds
			env.EnvVar & {CGO_ENABLED: "0"},

			// a globally consistent workdir
			env.Workdir & {path: "/work"},
		]
	}

	// a container after the code has built
	built: env.#Container & {
		@env()
		from: base
		steps: [
			// here we are passing the content directly as a string
			// load and add directories from your filesystem with #HostDir
			env.File & {path: "main.go", content: _goSrc},

			// run the build
			env.Exec & {args: ["go", "build", "-ldflags", "-w -s", "-o", "server", "main.go"]},
		]
	}

	// actual binary file for the server
	binary: env.#File & {@env(), path: "server", source: built}

	// start from an alpine, add the binary
	runner: env.#Container & {@env()
		from: "alpine:latest"
		steps: [
			// Sh is a wrapper around Exec to `sh -c <script>`
			env.Sh & {script: "apk add --update --no-cache ca-certificates"},

			// add the server binary
			env.File & {path: "/usr/bin/server", content: binary},

			// set the entrypoint
			env.Entrypoint & {args: ["/usr/bin/server"]},
		]
	}

	service: env.#Service & {@env(), name: "hello-service", ports:[{ port: 8080 }], source: runner}

	// caches we mount to the Go toolchain for across session caching
	caches: {
		build: env.Mount & { path: "/cache/go", source: env.#Cache & {name: "go-build-\(#ver)-\(#arch)"}}
		mods: env.Mount & {path: "/go", source: env.#Cache & {name: "go-mods-\(#ver)-\(#arch)"}}
	}
}

// you can @embed() files in CUE or load them with veg/env using #HostDir
_goSrc: """
	package main
	
	import (
	    "fmt"
	    "net/http"
	)
	
	func main() {
	    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	        fmt.Fprintf(w, "Hello, you've requested: %s\\n", r.URL.Path)
	    })
	
	    http.ListenAndServe(":8080", nil)
	}
	"""

```

### Build and Run a Container (with the catalog)

`hof env run veg-dev`

```cue
package veg

import (
	"github.com/hofstadter-io/hof/catalogs/env/packs"
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

dev: env.#Container & {
  @env(veg-dev)

  from: bases.debian13.default

  steps: [
    // customization
    packs.tool.zsh.customize,

    // deps for go/node/python -> c/c++ situations (like CGO)
    utils.apt.install & {#pkgs: ["gcc", "libc6-dev"]},

    // binary tool
    hof.File.linux,
    packs.tool.github.cli,

    // setup languages
    packs.lang.go.defaultSteps,
    packs.lang.cue.default,
    packs.lang.node.default,
    packs.lang.python.default,
    packs.lang.python.dev, // depends on node

    // tools for agents
    packs.tool.agents.lsp2mcp,

    // devops stuff
    packs.tool.hashicorp.terraform,
    packs.tool.hashicorp.packer,
    packs.tool.k8s.kubectl,
    packs.tool.k8s.helm,
    packs.tool.k8s.crane,

    // bind lsp servers, started on demand
    env.BindService & {service: packs.lang.go.lsp},
    env.BindService & {service: packs.lang.cue.lsp},
    env.BindService & {service: packs.lang.node.lsp},
    env.BindService & {service: packs.lang.python.lsp},
  ]
}
```

### Run an example

```sh
# get the repo
git clone https://github.com/hofstadte-io/hof && cd hof
make registry.start # todo, hof should handle this like it does for formatters
make dagger.start

# pick an example
cd examples/env/...

# or start something new
mkdir potato && cd potato
hof mod init vegg.ie/potato
hof mod get github.com/hofstadter-io/hof

# poke around
hof env list
hof env list ['^name$'] [-K '^kind$'] [-P '^path$'] [-S name|kind|path]

# sync, evaluates the DAGs, but doesn't export or expose services
hof env sync
hof env sync [...targets] [...flags]

# export a release bundle (from the root of hof)
hof env export -P dist -T my-test-rel

# veg out in a dev container (from repo root)
hof env run veg-dev
# > ls /usr/local/bin

# launch a full atprotocol network, including the permissioned pds
cd examples/env/atproto && hof env up

# make your own commands and flags, designed for your workflows
hof env [init, test, lint, ci, publish, deploy, ...]
hof env ... -t env=stg -t stack=app -t branch=main

# one-liners
hof env sync veg-dev
hof env run  veg-run
cd adk       && hof env test lint scan  # or `hof env ci`
cd atproto   && hof env up
cd inception && hof env turtles
```

<!-- # is k8s in dagger still a one-liner?
cd k8s && \
   hof env init && \
   hof env boot && \
   hof env lgtm && \
   hof env helm && \
   hof env test -->


### Code Organization

#### CUE (user) files

- [schemas/env](../../schemas/env)
- [catalogs/env](../../catalogs/env) (important, these are ordered to avoid cycles)
	- _note_, catalogs/env is an example and you can make your own, the only requirement is using the schemas/env
	- `utils/` helpers that only import veg/schemas/...
	- `bases/` operating systems and other basie images manually crafted, even from scratch
	- `packs/` abstractions, collections, and other reusable packs of env stuff
- [`.veg/`](../../.veg) the kitchen sink, it's all the things for this repo
- [examples/env](../../examples/env/)
	- `basic/` multi-stage build and 3-tier app
	- `advanced/` demonstrates more complex patterns and day-2 features
	- `adk/` commands example with matrix test and lint
	- `inception/` using docker, dagger, hof, helm, kubernetes from inside an env

Links to more examples can be found below.

#### Go (impl) files:

Everything is under `lib/env/...`

> [!TIP]
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

"points" are the targets, and they get sync'd, making it more a pull model than a push model
- so for ci, just decide what you want in the end, and let the dag engine figure out what actually needs to build and happen

Guidance on starting out with the commands

- `veg/env` tries to be context aware
  - knows the project or subdirectory you are in
  - knows the object type when handling args (veg-dist as example)
- the commands use the same flags and filtering
  - start with list, then run and a workful command
  - start with build or export before publish
- use `env.Terminal` or `-F/-N` to go interactive
- import `schemas/env/rrr:env` instead of `schemas/env` for stricter schemas. It's slower because of the recursive disjunctions, but it will catch errors earlier and narrow down messages, especially with the aliasv2 experiment we are relying on.

> [!TIP]
> `veg env` without any subcommands will run your custom commands

#### Help Text

_remember, hof -> veg, so veg ~> hof_

```sh
$ veg env --help
build, run, ship, and deploy environments (image, service, stack)

'veg env' looks for custom commands and treats them equally to builtin commands.
All commands have well known behaviors, depending on the $kind of a target.
'veg env list' and 'hof env sync' work with all '$kind's.
Most commands only work with a subset that makes sense or is explicit.
See their help text to learn more.

## Examples

# [...targets], or "points", are selected via args and flags
# list allows you to explore that space without syncing or triggering evaluation
veg env list|info ['^name$'] [-K '^kind$'] [-P '^path$'] [-S name|kind|path]

# sync, evaluates the DAGs, but doesn't export or make external alterations
# it can act as "real dry run" compared to the --dry-run flag for other commands
# without any args or flags, sync acts as a big test as well
veg env sync [...targets] [...flags]

# run, creates an interactive session and binds and dependent services
# this is closest to docker run or kubectl exec
veg env run [...target] [...flags]

# run, launches an services or stacks, similar to compose and helm
veg env up [...target] [...flags]

# export artifacts, local or remote, object storage and registries
veg env export -P release -T v0.4.3 -t dest=./release

# make your own commands and flags, designed for your workflows
# this is closest to Makefiles or package.json scripts
# define similar commands with the power of CUE and Dagger
veg env [init, test, lint, ci, publish, deploy, ...]
veg env ... -t env=stg -t stack=app -t branch=main

## Important References

./schemas/env    # the CUE schemas for what you can do in veg/env
./catalogs/env   # reusable CUE for all sorts of things from small to big
./examples/env   # simple and complex examples for you to play and fork

Usage:
  veg env [...target] [% ...cue] [flags]
  veg env [command]

Available Commands:
  export      export target points from an environment to outside world
  info        get details for target points in an environments
  list        list points in an environment
  run         run target point in an environment
  sync        sync target points in an environment
  up          starts target points in an environment

Flags:
      --all                    show all env targets, not just @env() ones
      --env-all                pass os.Env (everything)
      --env-file stringArray   path to a file with ENV vars to pass
      --env-var stringArray    key=value ENV vars to pass
      --fail-fast              fail at first error instead of attempting all targets
  -h, --help                   help for env
  -K, --kind stringArray       kinds to include, defaults to all
  -Z, --no-cache               bust the cache and force evaluation
  -N, --no-exit                Leave the TUI open after finishing
  -F, --on-failure             on failure, enter an interactive terminal, requires a tty
  -P, --parallel int           number of args or objects to process at once, they may be highly parallel internally (default 1)
  -R, --renderer string        output format [auto, plain, tty, dots, report (for ai)] (default "auto")
      --shh-file stringArray   path to a file with secret ENV vars to pass
      --shh-var stringArray    key=value secret ENV vars to pass
  -S, --sort stringArray       sort columns, can be used multiple times
      --unsafe                 set insecure root capabilities and privileged nesting, use at your own risk, needed for inception

Global Flags:
      --dry-run                  dry run certain commands
  -e, --expression stringArray   evaluate these CUE expressions only
      ...
```

Basically, the way this works is

1. CUE Value lattice (DAG.1) defines the Dagger OCI graph (DAG.2)
1. The CLI args/flags determine the point(s) you want to evaluate
1. To handle the request, for each point
    1. get the CUE Value, figure out the $kind
    1. recursively walk the CUE Value, build a Dagger AST, do some memoization
    1. Dagger Sync and perform which ever action the command is supposed to do

## Steps and #Stuff

Keep [schema/env](../../schemas/env/schema.md)
and [catalog/env](../../catalog/env)
handy for the details of the following.

Generally speaking...

- There are several groups or classes of statements, all prefixed by the `schemas/env.*` package identifier.
  - `env.Step` maps onto `With<Step>` and `Without<Step>` and can appear in `#Container: steps: [...]`
  - `env.#Stuff` maps onto resources, artifacts, and Dagger types. They are inputs, intermediates, outputs, or runnable.
  - `env.$Func` maps from one resource to one of the same or another, some `env.#Stuff` do some of these naturally too.
- It's a one way trip from CUE -> Dagger, you cannot for instance, use a directory listing or http response in CUE
  - `hof/flow` exists for this use case and some merging is on the roadmap.
  - The key requirement to maintain distinct operation modes, hermetic and yolo, with control over where, when, and how.

> [!IMPORTANT]
> We have swapped the semantics to `Token` | `#Token` from `WithToken()` and `Token()` with Dagger.
>
> 1. veg: `Dir` ~ dag: `WithDir()`
> 2. veg: `#Dir` ~ dag: `Dir()`

### Steps:

```
Dir               add a directory
File              add a file

Changes           applies #Changes to a container
Patch             applies a git-like patch or #Changes
PatchFile         applies a #PatchFile
RootFS            set the root FS of a container

EnvVars           add a set of env vars
SecretVars        add a set of secret vars
EnvFile           add an env var file
SecretFile        add a secret var file

Sync              force evaluation of the dagger graph
Exec              run any command as a container layer
Sh, Bash, Zsh     exec wrappers with a 'script' param
User              set the current user
Workdir           set the current workdir

Entrypoint        set container entrypoint
DefaultArgs       set container default args
DefaultTerm       set the terminal dagger uses when needed
Terminal          start a terminal at any or many point(s), directory or container

Temp              a temp volume for the next exec
Mount             mount a cache, file, directory, secret
UnixSocket        mount a unix socket at a path
BindService       bind another service to the container  (hint, dep graph)
Expose            mark a port for servin

                  these all remove from a container as a new layer
WithoutDefaultArgs
WithoutDirectory
WithoutEntrypoint
WithoutEnvVariable
WithoutExposedPort
WithoutFile
WithoutFiles
WithoutLabel
WithoutMount
WithoutRegistryAuth
WithoutSecretVariable
WithoutUnixSocket
WithoutUser
WithoutWorkdir
```
<!-- More to come...

- VsCode (like terminal, combo of them too)
- Chown
- ?Merge (not overwrite, doesn't exist yet)
- #DirToGit         git from a dir

Vscode            (we'll add a Step to open in vscode, or make something that does both, configurablely) -->


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
#RootFS           get the root FS for a container

#Changes          calculate the diff between two directories
#PatchFile        convert #Changes into a patch #File
#Shouldi          condition evaluation based on git diffs

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
#HostExec         run a command on the host (!outside of containers!)

#ExportDir        to host path
#ExportFile       to host path
#ExportImageFile  to a tarball
#ExportImage      to local engine
#PublishImage     to a registry


# helpers, not #things, but work / make them
#Flatten          create a new container from scratch and an existing container or directory

# many data formats available
#CuefigSBOM       returns a #File with the CUE/Dagger representation for any #Thing
```

## Examples

This repo is an excellent example, [`.veg/`](../../.veg) is a kitchen sink, it's all the things for this repo

[verdverm/testnet](https://github.com/verdverm/testnet) has a docker-compose like setup for [ATProtocol](https://atproto.com), build social apps, develop algos, run independent networks locally or for testing.

### Debian Base Container with Apt Caches

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

_from the root of this repo_

```sh
$ hof env list -e dist -S kind -S name
NAME               KIND          PATH                        EXTRA                                     
sbom-cuemod        cuefigSBOM    dist.sbom.cuemod            cuemod.cue                                 
sbom-fmt-black     cuefigSBOM    dist.sbom."fmt-black"       fmt-black.cue                              
sbom-fmt-prettier  cuefigSBOM    dist.sbom."fmt-prettier"    fmt-prettier.cue                           
sbom-github        cuefigSBOM    dist.sbom.github            github.cue                                 
sbom-veg-dev       cuefigSBOM    dist.sbom.dev               veg-dev.cue                                
sbom-veg-hof       cuefigSBOM    dist.sbom.hof               veg-hof.cue                                
sbom-veg-min       cuefigSBOM    dist.sbom.min               veg-min.cue                                
sbom-veg-ops       cuefigSBOM    dist.sbom.ops               veg-ops.cue                                
sbom-vscode        cuefigSBOM    dist.sbom.vscode            vscode.cue                                 
dist-cuemod        exportDir     dist.cuemod                 -> dist/cuemod                             
dist-github        exportDir     dist.github                 -> dist/github                             
dist-meta          exportDir     dist.meta                   -> dist/meta                               
dist-vscode        exportDir     dist.vscode                 -> dist/vscode                             
dist-fmt-black     publishImage  dist.images."fmt-black"     -> host.docker.internal:5000/fmt-black     
dist-fmt-prettier  publishImage  dist.images."fmt-prettier"  -> host.docker.internal:5000/fmt-prettier  
dist-veg-dev       publishImage  dist.images.dev             -> host.docker.internal:5000/veg-dev       
dist-veg-hof       publishImage  dist.images.hof             -> host.docker.internal:5000/veg-hof       
dist-veg-min       publishImage  dist.images.min             -> host.docker.internal:5000/veg-min       
dist-veg-ops       publishImage  dist.images.ops             -> host.docker.internal:5000/veg-ops       
```

```sh
# validate a release is ready to go
$ hof env sync -e dist -T v0.7.0-alpha.2

# export and publish artifacts for real
$ hof env export -e dist -T v0.7.0-alpha.2 -t reg=ghcr.io/hofstadter-io
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
  from: bases.debian13.minimal
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

Reasons to have CUE in your toolbox

1. It's an approachable, logical language. It's a can be a good thinking tool like Haskell if it's not a daily driver for you. It has mathsex appeal.
2. It deals with config and data, rosetta stone, not just for formats. JQ on steroids
3. It can be surgical, or global if you prefer. Add schema, transform data, or enforce policy
4. Everything is pretty much ETL anyway, CUE is the glue or meta layer (#devops)
