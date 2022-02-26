# Ideas

- can we build `hof/flow` in such a way that
  Dagger can build their tool on top of it?
- `hof/ui` command
- Weekly CUE updates (edited video, newsletter, social posts)

## Tonight

Everynight should include...

- [ ] docs

---

Other:

- [ ] use decoder structs
  - [ ] use in simple task & middleware first
    ts: getenv, watch, readfile
    ms: print
  - [ ] when revamping api.Call
- [x] ensure `@print()` is working (pre/post?)
- [ ] youtube videos

Do all these together, per task or similar
  - [ ] decode type? (consider for some)
  - [ ] docs
  - [ ] tests
  - [ ] task schemas
  - [ ] cleanup existing tasks

## Questions

should middleware and/or tasks use $<key>
- $log: warn: "my message"
  $log: warn: [string]: _
  $log: global: ...? (level, format)
  $log: local: ...?  (level, ...?)
- $task: ??? (probably not)
- can we support one, other, or both with config / flags
- should they imply different things? what would that be?
  - task vs middleware config
  - global / local config / vals

Using `$` as a signifier means we can use real values and the evaluator to check correctness.
`@` is OK for things that don't need more than minimal configuration, typically a name like `@flow` and `@task`

middleware which is scoped both global and task local, eventually grouped tasks (when nested tasks is supported)

We should probably try to get both UX versions in the next release (0.6.2) so we can test and get feedback

### v0.6.2

`hof/flow` command

- [x] CLI sugar:
  - [x] run solo named flows
  - [x] easier flows with: @flow/name
  - [x] easier tags with: +tag=value
- [ ] first pass at docs
- [ ] tests / examples
- [ ] task schemas for import
- [ ] more tasks
  - [x] @pool for controlling max parallelism (middleware technically)

- [ ] fix / finish tasks & anything else
  - [ ] API call response enrichment (where else?)
  - [x] remove global vars

- [x] centralizations (middleware/basetask)
  - [x] progress / stats
  - [x] finalized printing of stats

### v0.6.3

CLI sugar:
  - @flow/name
  - +tag=value
  - args for pipelines
   - use -- for everything after
   - populate `$/@args: ` struct

- better nested context (path, stats)
- $syntax (see other notes for details)
- logging
  - [ ] log level (unify with global needs and share)
- structs to decode tasks / middleware with (simpler hopefully)

### Frame out hof/flow API

- [x] deglobalize TaskRegistry & Context
- [x] Initial Task Value - Root & Path, can we remove the local reference to a complete task and infer the need to run (cue/issue #1535)
- [x] base task default type (composite) + interface(s)
- [x] task IDs
- [ ] print task dependency tree
  - [x] simple list 
  - [x] simple with stats
  - [ ] tier/topo sorted
- [ ] middleware
  - [ ] wait group
  - [ ] stats / prometheus / metrics
  - [ ] progress / open telemetry
  - [ ] bookkeeping / history
  - [ ] finalize stage (i.e. write final history)
- [ ] run flow without `-f` flag
- [ ] history saving (as middleware)

want reusable and extendable

- remove task.Final
- collect stats
- abstract (or custom/extended) context, tasks may still need to be tied to it?
- [x] remove global taskmaker
- make task diagram when done -> as built
  - mermaid (or other) diagram generated
- middleware...
- well known $id like fields over attributes?
- eventing system? (outside of tasks, but about the process)
  (both)

---

Plugin is a reusable concept made of
- middleware
- custom tasks
- new attributes
- context / base task extension
- event emiters / handlers

Eventually want to extract builtin default stuff
to plugins

---

hof adhoc generator for task implementations
(lots of repeat going on)

---

- [ ] task middleware ?!
    - local (pull) | global (push)
    - how to register, code only
    - two funcs to define
      - discovery
      - running
    - how to expose config, per/all flow granularity
    - (These are Go API questions)

Code:

- [ ] wait group (do we need this?)
- [ ] sync.Pool for controlling max parallelism (global setting, local pools)
- [ ] task specific "loop" support
- [ ] API call response enrichment (where else?)
- [ ] review tasks to find more fix/finish

Docs: 

- [ ] first pass at docs

Prior art:

- [ ] make list for task engines / frameworks
- [ ] features / pros / cons / gaps

---

- [ ] (at least quit & main) select, random like Go, or also preferred order?
- [ ] (don't need, sync pool is what we were after) "go routine" (really tasks are already if they do not exit)
- [ ] broadcast, chan 1->Many, this probably needs a different type in the Context chanmap
- [ ] cron / time chan producer, which wraps another task
- [ ] write youtube video data
- [ ] oauth token refresh
- [ ] set playlists declaratively
- [ ] refactor twitch flows



### v0.7.0

after CUE v0.4.3

- [ ] support dynamic tab completion of hof/{flow,gen,dm}
- [ ] support other data formats for tags (int,bool,[]string) 
- [ ] can we eliminate our patched cue? (MakeError & task.Final)
- [ ] logging middleware (move to separate concept?)

- [ ] deploy flow(s) that can be embedded as a top-level command
- [ ] consider rewrite of discovery / running engine
      (this might really be a requirement)
- [ ] required helper so we can use to require fields
- [ ] validate tasks with schema
- [ ] playwright for more oauth automation
- [ ] ( not just csp, API, kafka, irc, websocket )
- [ ] easier to customize / extend / sandbox tasks & middleware

- [ ] history (persistent), get results from past runs, diff view
     - official flow to import and snyc to your favorite cloud

- [ ] required fields (from schema) as a field & check in tasks
      should/can we check a subtask like handler upfront?
      what if the task conditionally depends on 

- [ ] Per task loop support, 
      - task needs to support this, hard to generalize for a first pass, really code dependent?
      - but maybe not, don't spend much time trying to figure this out for a first pass
---

### Flows & Tasks

`hof` work

- [x] migrate to hof/flow cmd
- [ ] hof/lib/st.Required
- [ ] task schemas (github.com/hof-io/hof/flow/tasks/...)
- [ ] *** internal validation
- [ ] *** task types for decoding (partially?)
- [ ] update to flow based
  - [ ] structural txtar tests
- [ ] improve test situation in github actions

```
@flow(init)
@use()
```


- [ ] new 'flows' repo in hof github as a
  - [ ] docker containers
    - [ ] task library for importing
    - [ ] start / stop

- [ ] implement oauth workflow for youtube
- [ ] declarative youtube videos & playlists like kubernetes
- [ ] update twitch library from learnings with youtube


### hof/flow examples

- [ ] schemas for builtin tasks
- [ ] make this list
  - [ ] examples for available tasks
  - [ ] composite tasks like docker
  - [ ] links to tools using flow (streamer-tools)

- save all IRC messages to DB
- bookmarks and hn upvotes
- change my lights
- replace helm (need native topo sort)
- OAuth workflow

### hof/flow internals

- [ ] metrics and progress
  - [ ] TaskStats: struct and centrally collecting
  - [ ] Print progress (and stats) by flag
  - [ ] failure mode for tasks, some common schema across all tasks? @onfail()

- [ ] i/o centralization
  - [ ] debug/verbose flag to print tasks which run
  - [ ] stats for tasks & pipelines, chan to central
  - [ ] obfescate secrets, centralized printing (ensure that is the case in every task / run)

- [ ] exec improvements
  - [ ] many options not enabled yet
  - [ ] for exec: pipe stdin/out/err to files
  - [ ] some way to run in background, and then kill / exit later?

- [ ] async / client listener
  - [ ] chan / mailbox
  - [ ] waitgroup / mutex?
  - [ ] kill chan, also need to catch signals in main(?) and pass down / do right thing
    - [ ] how to tell (server / bg exec'd process) to stop (oauth localhost as example)
  - [ ] websockets

- something that loops over input list and produces messages / tasks (?)

- [ ] sql
  - [x] sqlite
  - [ ] postgres
  - [ ] mysql

- [ ] msg
  - [ ] rabbitmq
  - [ ] kafka
  - [ ] nats

- [ ] k/v
  - [ ] redis
  - [ ] memcache
  - [ ] gcs
  - [ ] s3

- [ ] fs
  - [x] watch
  - [ ] load into mem

- [ ] obj
  - [ ] elasticsearch
  - [ ] mongo

- [ ] mouse/keyboard automation
  - [ ] Browser - https://github.com/playwright-community/playwright-go
  - [ ] OS level - https://github.com/go-vgo/robotgo

- [ ] hof gen as a task

- [ ] server
  - [ ] logging levels

Other: 

- specify CLI args to flow for command to run (avoid -p)
- better (boolean,regexp) logic for selecting pipeline(s)
- temp files / dirs
- command line prompt
- support for fs.FS (https://github.com/hack-pad/hackpadfs)
  - could abstract away S3/GS

### Build other things cuetils/run

### More todo, always...

Exec & Serve & async

- [ ] write directly to file for stdio, is it a concrete string?
- [ ] something like a goroutine, similar to api.Serve/pipeline
- [ ] message passing, via chans, websockets, kafka/rabbit

Bugs?

- [ ] prevent exit when error in handler pipelines?
- [ ] rare & racey structural cycle
- [ ] cuetils flow args for CWD all behave differently
  - [ ] `<no args>` -> no input
  - [ ] `*.cue` -> refs across files not found
  - [ ] `./` -> works
- [ ] imported flows that have typos don't throw error, rather they don't show up silently
- [ ] in @task(ext.CueFormat) [ unsupported type <nil>...? ]

Helpers:

- canonicalize (sort fields recursively)
- toposort

List processing:

- jsonl
- yaml with `---`
- CUE got streaming json/yaml support
- if extracted value is a list?

