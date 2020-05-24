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
  setup           α     setup the hof tool
  init            α     create an empty workspace or initialize an existing directory to one
  clone           α     clone a workspace or repository into a new directory

Model your designs, generate implementation, run anything:
  modelset        Ø     create, view, diff, calculate / migrate, and manage your data models
  gen             ✓     generate code, data, and config from your data models and designs
  run             α     run polyglot command and scripts seamlessly across runtimes
  runtimes        α     work with runtimes (go, js, py, bash, custom)

Learn more about hof and the _ you can do:
  doc             Ø     Generate and view documentation.
  tour            Ø     take a tour of the hof tool
  tutorial        Ø     tutorials to help you learn hof right in hof

Download modules, add content, and execute commands:
  mod             β     mod subcmd is a polyglot dependency management tool based on go mods
  add             α     add dependencies and new components to the current module or workspace
  cmd             α     run commands from the scripting layer and your _tool.cue files

Manage resources (see also 'hof topic resources'):
  info            α     print information about known resources
  label           α     manage labels for resources and more
  create          α     create resources
  apply           α     apply resource configuration
  get             α     find and display resources
  edit            α     edit resources
  delete          α     delete resources

Configure, Unify, Execute (see also https://cuelang.org):
  (also a whole bunch of other awesome things)
  def             α     print consolidated definitions
  eval            α     print consolidated definitions
  export          α     export your data model to various formats
  fmt             α     formats code and files
  import          α     convert other formats and systems to hofland
  trim            α     cleanup code, configuration, and more
  vet             α     validate data
  st              α     recursive diff, merge, mask, pick, and query helpers for Cue

Manage logins, config, secrets, and context:
  auth            Ø     authentication subcommands
  config          β     manage local configurations
  secret          β     manage local secrets
  context         α     get, set, and use contexts

Examine workpsace history and state:
  status          α     show workspace information and status
  log             α     show workspace logs and history
  diff            α     show the difference between workspace versions
  bisect          α     use binary search to find the commit that introduced a bug

Grow, mark, and tweak your shared history (see also 'hof topic changesets'):
  include         α     include changes into the changeset
  branch          α     list, create, or delete branches
  checkout        α     switch branches or restore working tree files
  commit          α     record changes to the repository
  merge           α     join two or more development histories together
  rebase          α     reapply commits on top of another base tip
  reset           α     reset current HEAD to the specified state
  tag             α     create, list, delete or verify a tag object signed with GPG

Colloaborate (see also 'hof topic collaborate'):
  fetch           α     download objects and refs from another repository
  pull            α     fetch from and integrate with another repository or a local branch
  push            α     update remote refs along with associated objects
  propose         α     propose to incorporate your changeset in a repository
  publish         α     publish a tagged version to a repository
  remotes         α     manage remote repositories

Local development commands:
  reproduce       Ø     Record, share, and replay reproducible environments and processes
  jump            α     Jumps help you do things with fewer keystrokes.
  ui              Ø     Run hof's local web ui
  tui             Ø     Run hof's terminal ui
  repl            Ø     Run hof's local REPL
  pprof                 go pprof by setting HOF_CPU_PROFILE="hof-cpu.prof" hof <cmd>


Send us feedback or say hello:
  feedback        Ø     send feedback, bug reports, or any message :]
                        you can also chat with us on https://gitter.im/hofstadter-io

Additional commands:
  help                  help about any command
  topic                 additional information for various subjects and concepts
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal

Additional topics:
  schema, codegen, modeling, mirgrations
  resources, labels, context, querying
  workflow, changesets, collaboration

(✓) command is generally available
(β) command is beta and ready for testing
(α) command is alpha and under developmenr
(Ø) command is null and yet to be implemented

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

Use "hof [command] --help / -h" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.
```
