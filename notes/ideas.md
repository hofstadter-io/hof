# Ideas


### Tests...

txtar in 

### Docs...

- [ ] start hof/docs folder(s)
- [ ] task schemas (github.com/hof-io/hof/flow/tasks/...)
- [ ] task reference (autogen the majority)


### Other

Go funcs:

- rename currenct to `*Globs`
- pure Go implementations
- funcs that take values

CLI:

- Support expression on globs, to select out a field on each file
- move implementation?

### futurology

- @label(), but also part of evaluation? (available for gens and flow)
- diff lists, @id(), how to detect renames and position changes and optimize?


## upstream & references

#### Memory issues

(we have not seen this yet with the twitch IRC bot which had lots of activity)

https://github.com/lipence/cue/commit/6ed69100ebd62509577826657536172ab46cf257

#### cue/flow

return final value: https://github.com/cue-lang/cue/pull/1390
