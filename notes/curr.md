# curr

### Docs, Demos, & Mods

- [ ] docs
	- show most minimal schema for hof dm & gen
	- inline partials:
- [ ] record demo
	- [ ] as a user creating an app only, import gens
  - [ ] adhoc to module, show new mode & flags
  - [ ] generating types across languages (also as example for docs)
- [ ] Mods
  - [ ] `hofmod-openapi` need to digitize notes
  - [ ] `hofmod-cli` update & enhance
	- [ ] `hofmod-type` build out a first version
  - [ ] `hofmod-app` from `hofmod-{server,openapi,...}`
	  - show how to order subgens

```
{{ range $items }}
{{ template "inline-partial" . }}
{{ end }}

{{ define "inline-partial" }}
...
{{ end }}
```



### next

bugs:

- loading only yaml broken (need repro case)


v0.6.5

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

v0.6.x

- [ ] real datamodel upgrades & efforts
- [ ] hof faux
- [ ] add cue template helper & datafile
- [ ] override templates & partials, propagate through subgens (layer back up)
- [ ] Gen (top-level) commands (exec), gens can provide some nice defaults

