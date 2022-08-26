# curr


# today

- [x] get templates rendering from creator gen
- [x] render CUE files (generators)

Path issues to resolve

- [ ] potential template lookup issues for all other things here
- [ ] nested gen (don't want to force top-level creators)
- [ ] local dev for creator mods and ../ like paths

Two mods to test with:

- [ ] hof create github.com/hofstadter-io/hofmod-cli   (root gen)
- [ ] hof create github.com/hofstadter-io/hofmod-types (nested gen)

We support creators in subdirs.
We probably need to walk up remote path to find git repo like Go.
Use this to prototype revamp of `hof mod`




### Docs, Demos, & Mods

- [ ] docs
	- show most minimal schema for hof dm & gen
	- inline partials:
- [ ] record demo
	- [ ] as a user creating an app only, import gens
  - [ ] adhoc to module, show new mode & flags
  - [ ] generating types across languages (also as example for docs)
- [ ] Mods
	- [ ] `hofmod-types` build out a first version
  - [ ] `hofmod-openapi` need to digitize notes
  - [ ] `hofmod-cli` update & enhance
  - [ ] `hofmod-app` from `hofmod-{server,openapi,...}`
	  - show how to order subgens

Document:

```
{{ range $items }}
{{ template "inline-partial" . }}
{{ end }}

{{ define "inline-partial" }}
...
{{ end }}
```

### todo...


Changes to shadow dir usage:

- [ ] ensure if cue.mod is not found, we make .hof in the local dir
- [ ] tests for all the various situations
- [ ] how does this work for nested generators?
  - what if a single generator uses the same subgen twice?


Formatting:

- configured via generator / env vars?
- show how to do in several CI systems? (at least GHA)




### next

bugs:

- loading only yaml broken (need repro case)
- [ ] `hof fmt file.cue` is not working
- [ ] `hof fmt start all[` segfaults

v0.6.x

- [ ] hofmod-cli revamp
- [ ] hofmod-openapi & formatters
- [ ] hof create with https://github.com/AlecAivazis/survey
- [ ] hof mod iteration
  - [ ] get
	- [ ] move cue.mods -> cue.mod/hof-{mod,sum}.cue
	- [ ] symlink replaces


v0.6.x



- [ ] generator info (help to know what you can do, what you can override, docs)
- [ ] integrate datamodel history to hof gen
- [ ] hof create ... (also from github repos) easier startup for gen user only mode
  - [ ] generators should have create templates as well
	- [ ] use -t (tag) to inject values? (or something else to not conflict with cue?)
- [ ] hof st (import bulk commands from cuetils)
- [ ] hof mod get & tidy, move to cue.mod/hof.cue
- [ ] hof doc
- [ ] remove deps
	- [ ] github.com/aymerick/raymond (not really used)
	- [ ] github.com/bmatcuk/doublestar/v4 ?

- [ ] real datamodel upgrades & efforts
- [ ] hof faux
- [ ] add cue template helper & datafile
- [ ] override templates & partials, propagate through subgens (layer back up)
- [ ] Gen (top-level) commands (exec), gens can provide some nice defaults

