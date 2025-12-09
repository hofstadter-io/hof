# ./lib/agents


## Architecture


Core features

- one server, main clients, each with many sessions that can run while they are away
- both REST and websocket connections are used
  - share message types and format
  - websocket allows for always on async streaming and broadcasting
  - REST allows for synchronous requests or processing
- ADK manages events, sessions, and agents
- `Runtime` > `Client` > `Session` > `Env`
  - `Runtime` holds the core services [ADK,DB,Dagger,Server,Clients]
  - `Client`
- `

### Runtime & Server - REST & Websocket server

Important files:

- [core runtime and interfaces](lib/agent/runtime/*.go) - messy :[
- [services for ADK]()

- `github.com/labstack/echo` server
- `github.com/google/adk-go` framework


### Services

- sessions - agents, events, and history
- dagger - filesystems and environments
- artifacts - key/blob store for files
- memories - searchable memories (not sure this is needed?)
  - this should just be a tool

### Virutalized Environments - connects VS Code concepts to Dagger concepts

- file system and terminal based on `dagger.Container()`
- expose through REST & Websocket messages




### Repository Layout

```sh
.
├── agents
│   └── cue.go (maps CUE -> Go types)
├── AGENTS.md  (this file)
├── cmd
│   └── run.go (old, meant to run an agent)
├── extension
│   └── cmd.go (current, background server entrypoint for the VS Code Extension)
├── models
│   └── gemini.go (boring)
├── runtime (core, a bit messy)
│   ├── api_fs.go
│   ├── client.go
│   ├── dagger
│   │   └── client.go (singleton client)
│   ├── handlers/
│   │   └── ws/
│   ├── run.go
│   ├── runtime.go
│   ├── session.go
│   └── services (next, we need to refactor)
│       ├── artifact/
│       ├── container.go (curr, this is what we are working on)
│       ├── dagger/ (curr, this is what we are working on)
│       ├── db.go
│       └── session/
└── tools
    ├── filesys/... (old, boring)
    ├── mathy.go (old, boring)
    ├── mcp (current, mcp wrappers for agents)
    │   ├── github.go
    │   ├── local.go
    │   └── tavily.go
    └── meta (current, tools for agents)
        ├── cache.go
        ├── exec.go
        └── filesys.go
```
