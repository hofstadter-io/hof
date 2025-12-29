# hof/env

`hof/env` is a CUE interface to something like Docker + Compose

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
  -h, --help              help for env
  -Z, --no-cache          bust the cache and force evaluation
  -N, --no-exit           Leave the TUI open after finishing
  -F, --on-failure        on failure, enter an interactive terminal, requires a tty
  -P, --progress string   output format [auto, plain, tty, dots, report (for ai)] (default "auto")
```