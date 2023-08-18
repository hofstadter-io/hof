---
title: Modules

weight: 75
---


{{<lead>}}
`hof mod` helps you manage
CUE modules and dependencies.
They are based on Go modules,
can be backed by git or oci,
support public and private visibility,
and make working with CUE modules easy!
{{</lead>}}


<div id="module-cast" class="asciinema"></div>

<br>
<br>


A CUE module is a collection of packages and
is recognized by a `cue.mod/module.cue` file
containing a `module: "github.com/org/repo"` field.
The modules and packages work like other languages,
allowing you to organize your code and later import it elsewhere.
CUE is already module aware but lacks dependency management.
`hof mod` fills that gap to make working with modules easy.

`hof mod` is based on Go modules, so you will recognize
both the commands and the dependency file.
In fact, `hof` reuses much of the algorithm.
The main difference is the file location and format.
Another notable difference is that `hof` has no
global servers or registries, and instead talks
directly to the code host to avoid adding more
middlemen in your software supply chain.

A brief introduction is below while
the following pages go into detail.

{{<childpages>}}


### Creating a Module

To create a new module, run `hof mod init <module-path>`.
This will create a `cue.mod/module.cue`
If you have already manually made the
`cue.mod/module.cue` file, then you can skip this step.

{{<codeInner>}}
hof mod init github.com/hofstadter-io/hof
{{</codeInner>}}

It is best to use the same path as the Git repository.
This makes your code immediately reuseable by pushing commits or a semver tag.

{{<codeInner>}}
module: "github.com/hofstadter-io/hof"
cue: "0.6.0"

require:  { ... } // your dependencies
indirect: { ... } // dependency-dependencies
replace:  { ... } // replaces during development
{{</codeInner>}}


### Adding a Dependency

Use `hof mod get` to add and fetch a new dependency.
Hof recogonizes semver tags,
several special version names,
and allows for any commit hash.

{{<codeInner>}}
hof mod get github.com/hofstadter-io/hof@v0.6.8
hof mod get github.com/hofstadter-io/hof@v0.6.8-beta.6
hof mod get github.com/hofstadter-io/hof@latest   // latest semver
hof mod get github.com/hofstadter-io/hof@next     // next prerelease
hof mod get github.com/hofstadter-io/hof@main     // latest commit on branch
hof mod get github.com/hofstadter-io/hof@475328015adf6d102e5227a646e63f6a2b23119f
{{</codeInner>}}

### Updating Dependencies

Use `hof mod tidy` to inspect your imports,
update the dependency list, and fetch any missing modules.
This command will update both `module.cue` and `sums.cue`,
which is used to store module hashes.
You can verify code integrity later by running `hof mod verify`.

You can update dependency versions with `hof mod get`

- `hof mod get <module-path>@latest` or `@next` for a prerelease
- `hof mod get all@latest` to update all dependencies to their most recent version


### Fetching Dependencies


Hof supports two methods for fetching dependencies into your project.

1. `hof mod vendor` will copy the files into `cue.mod/pkg/...`
2. `hof mod link` will symlink the directories into `cue.mod/pkg/...`

The choice can depend on personal preference or system requirements.
(Windows does not support symlinks)
Linking makes co-development with local replaces easier,
as changes are automatic and 
do not require you to run `hof mod vendor` to get updates
from another directory.

<br>

<p><button class="btn btn-primary" type="button" data-bs-toggle="collapse" data-bs-target="#collapseExample" aria-expanded="false" aria-controls="collapseExample">
	hof mod help
</button></p>
<div class="collapse" id="collapseExample">
{{<codePane file="code/cmd-help/mod" lang="txt" title="$ hof mod help">}}
</div>


### Dependency Cache

`hof mod` fetches and copies module code into your system's cache directory.
This acts as a per-users, local module cache and helps reduce network requests and disk space.
A module will only be fetched once on a system.
You can find the cache directory in `hof version` and
you can delete it with `hof mod clean`.

Note, Hofstadter does not run any proxies or sumdb and all requests go directly to the module host.

