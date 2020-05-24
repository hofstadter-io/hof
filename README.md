# hof - a polyglot tool for building software

The `hof` tool started as a code generation framework
and has since expanded to helping make dev work fun again.
This is a large body of work in progress,
contributions are warmly welcomed.

### getting started

You will have to download `hof` the first time.
After that `hof` will prompt you to update and
install new releases as they become available

```text
# Install (Linux, Mac, Windows)
curl -LO https://github.com/hofstadter-io/hof/releases/download/v0.5.3/hof_0.5.3_$(uname)_$(uname -m)
mv hof_0.5.3_$(uname)_$(uname -m) /usr/local/bin/hof

# Shell Completions (bash, zsh, fish, power-shell)
echo ". <(hof completion bash)" >> $HOME/.profile
source $HOME/.profile

# Show the help text (also seen below)
hof

# Setup hof (optional)
hof setup
```

You can always find the latest version from the
[releases page]

Go Docs: https://pkg.go.dev/mod/github.com/hofstadter-io/hof

Chat: https://gitter.im/hofstadter-io


### top-level commands

```text
hof - a polyglot tool for building software

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]


Setup hof and create workspaces and datasets:
  setup                 Setup the hof tool
  init                  Create an empty workspace or initialize an existing directory to one
  clone                 Clone a workspace or repository into a new directory

Model your designs, generate implementation, run anything:
  modelset              create, view, migrate, and understand your modelsets.
  gen                   generate code, data, and config
  run                   run polyglot command and scripts seamlessly across runtimes
  runtimes              work with runtimes

Learn more about hof and the _ you can do:
  doc                   Generate and view documentation.
  tour                  Take a tour of the hof tool
  tutorial              Tutorials to help you learn hof right in hof

Download modules, add content, and execute commands:
  mod                   mod subcmd is a polyglot dependency management tool based on go mods
  add                   add dependencies and new components to the current module or workspace
  cmd                   Run commands from the scripting layer

Manage resources (see also 'hof topic resources'):
  info                  print information about known resources
  label                 manage labels for resources and more
  create                create resources
  apply                 apply resource configuration
  get                   find and display resources
  edit                  edit resources
  delete                delete resources

Configure, Unify, Execute (see also https://cuelang.org):
  (also a whole bunch of other awesome things)
  def                   print consolidated definitions
  eval                  print consolidated definitions
  export                export your data model to various formats
  fmt                   formats code and files
  import                convert other formats and systems to hofland
  trim                  cleanup code, configuration, and more
  vet                   validate data
  st                    Structural diff, merge, mask, pick, and query helpers for Cue

Manage logins, config, secrets, and context:
  auth                  authentication subcommands
  config                Manage local configurations
  secret                Manage local secrets
  context               Get, set, and use contexts

Examine workpsace history and state:
  status                Show workspace information and status
  log                   Show workspace logs and history
  diff                  Show the difference between workspace versions
  bisect                Use binary search to find the commit that introduced a bug

Grow, mark, and tweak your shared history (see also 'hof topic changesets'):
  include               Include changes into the changeset
  branch                List, create, or delete branches
  checkout              Switch branches or restore working tree files
  commit                Record changes to the repository
  merge                 Join two or more development histories together
  rebase                Reapply commits on top of another base tip
  reset                 Reset current HEAD to the specified state
  tag                   Create, list, delete or verify a tag object signed with GPG

Colloaborate (see also 'hof topic collaborate'):
  fetch                 Download objects and refs from another repository
  pull                  Fetch from and integrate with another repository or a local branch
  push                  Update remote refs along with associated objects
  propose               Propose to incorporate your changeset in a repository
  publish               Publish a tagged version to a repository
  remotes               Manage remote repositories

Local development commands:
  reproduce             Record, share, and replay reproducible environments and processes
  jump                  Jumps help you do things with fewer keystrokes.
  ui                    Run hof's local web ui
  tui                   Run hof's terminal ui
  repl                  Run hof's local REPL
  pprof                 Go pprof by setting HOF_CPU_PROFILE="hof-cpu.prof" hof <cmd>


Send us feedback or say hello:
  feedback              send feedback, bug reports, or any message :]
                        you can also chat with us on https://gitter.im/hofstadter-io

Additional commands:
  help                  Help about any command
  topic                 Additional information for various subjects and concepts
  update                Check for new versions and run self-updates
  version               Print detailed version information
  completion            Generate completion helpers for your terminal

Additional topics:
  schema, codegen, modeling, mirgrations
  resources, labels, context, querying
  workflow, changesets, collaboration

Flags:
      --account string               the account context to use during this hof execution
  -E, --all-errors                   print all available errors
      --billing string               the billing context to use during this hof execution
      --config string                Path to a hof configuration file
      --context string               The of an entry in the context file
      --context-file string          The path to a hof context file
      --global                       Operate using only the global config/secret context
  -h, --help                         help for hof
      --ignore                       proceed in the presence of errors
      --impersonate-account string   account to impersonate for this hof execution
  -l, --label strings                Labels for use across all commands
      --local                        Operate using only the local config/secret context
      --log-http string              used to help debug issues
  -p, --package string               the package context to use during this hof execution
      --project string               the project context to use during this hof execution
  -q, --quiet                        turn off output and assume defaults at prompts
      --resources-dir string         directory for discovering resources
      --runtimes-dir string          directory for discovering runtimes
      --secret string                The path to a hof secret file
  -S, --simplify                     simplify output
      --strict                       report errors for lossy mappings
      --trace                        trace cue computation
  -T, --trace-token string           used to help debug issues
      --tui                          run the command from the terminal ui
      --ui                           run the command from the web ui
  -v, --verbose string               set the verbosity of output
      --workspace string             the workspace context to use during this hof execution

Use "hof [command] --help" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.
```
