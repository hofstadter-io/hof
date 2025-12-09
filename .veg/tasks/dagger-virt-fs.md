## Virtualized Environments for VS Code, backed by Dagger

Overview:

- We have a VS Code extension and want to implement the `vscode.FileSystemProvider` interface
  - We have many of these in a VS Code workspace, the use a unique `eid` (environment id) for association
  - The extension code is elsewhere and talks to this server over an api and/or websocket as needed
- This server implements the API and Websocket server
  - Mirror the VS Code FS interface, share the same endpoints and req/resp schemas, Websocket wraps the equivalent REST body like this: `{ type: string, payload: any }`
- This server implements an Agentic framework with ADK
  - Environments are attached to Sessions
  - The Agent interacts with the Environment (file sys / terminal) through Tool calls `fs_read`, `fs_write`, `exec`, etc...
  - these are capture in events
- This server implements a DaggerService which provides a simpler interface for our VS Code and ADK needs
  - It needs to support situations in VS Code where agents are not involved
  - When agents are involved, we have to be aware of and send events two-way, the process is iterative with humans using VS Code and changes to files there need to be accumulated and added as part of the next user message. (As functionResp perhaps?)


Details:
- `Directory` based for core file system
- Reusable package mirroring the `vscode.FileSystemProvider` interface.
- Needs to support multiple sessions per user and snapshots per session (it is actually a tree, like git and branches, but with Dagger/BuildKit)

- VS Code Extension Server
  - REST and Websocket handlers parse inputs to call this package.
- ADK Framework
  - Service
  - Tools for Agents to use
- Types (req/resp/config) for each of the file system calls need to be aligned as much as possible.
  - 

I have prototyped the basic workings of this, but it is a mess.


paths to files with relevant or related code:

- [REST] `lib/agent/runtime/api_fs.go`
- [Tool] `lib/agent/tools/meta/filesys.go`

Problems:

1. Both the API and Tools
  - Interact with dagger
  - Have types defined for them
  - No websocket for the FS yet
  
2. The api mutations are not capture in a sessions
  - this needs to be conditional it being part of an agent session

3. I used `sid` for session in the prototype, but the concept is more general now
  - use `eid` for Environment ID going forward
  - this is used in
    - the `vscode.Uri` query params
    - in payloads of messages / requests / handlers
    - (todo) in the database to map eid -> dagger ID (see below)
  - Session State
    - currently holds `{ sid, origfs, dagger, origrv, runenv }`
    - need better names `{fs,env}-{orig,curr}` (consistent, not tied to a tech)
    - need short hash as value instead of dagger ID
    - "runenv" is now Environment

  
- types are in the tool, while we arent' working on this right now, they have the shape the LLM sees and the calls to Dagger that back them.

Ideally we can abstract the common payloads out from the various places I already have code, into a single, source-of-truth Go type

**Current Target**:

IMPORANT:  - the reusable virtualized interface

We may need


**Current Plan**:

We will work on this as we go iterate together

Goal: VS Code -> REST/Websocket -> DaggerService (no agent stuff yet)

Long-term plan:

1. Write the Dagger EnvironmentService 
  - files: `lib/agent/services/environment/dagger_fs_{query,mutate}.go`
  - notes: needs a singleton that both ADK and VS Code can use (independently, just VSC for now)
2. Write the Echo API to mirror VS Code FileSystemProvider
  - files: `lib/agent/runtime/handlers/api/vscode_filesystem.go`
  - notes: this is probably where we do the translation from VS Code -> Dagger
3. Write the VS Code FileSystemProvider that uses the API
  - files: `extensions/vscode/extension/src/services/filesystemProvider.ts`
  - notes:
    - a few of the query endpoints are already implemented
    - we can ignore watch()
    - we need to figure out notifications so the VS Code explorer auto refreshes
4. VS Code UI
  - explorer context menu contribution to add veg area
    - "open dir with veg", devs typically have an actual filesystem path open
    - they should be able to, right click and...
      - open the real folder into a virtual folder
      - open a real subdir as a virtual folder
      - do the same from a virtual folder
      - the result is a lineage tree of folders as workspace entries
        - need to record the parent, as a vscode.Uri with eid as the parameter

## **Current Task**:

1. Write the Dagger EnvironmentService
  1. Singleton Client
    1. CreateSingleton that takes `*gorm.DB`
    2. Get() that returns the our wrapped Client
  2. "SDK" for core functions we need for VS Code
    1. Attached to the Singleton Client so can access DB
  3. Parameterized by
    1. `eid` (the environment id)
  3. Mutations should return a new `eid` and save it do the datbase
  4. Gorm setup for a single table mapping `eid` -> `VARCHAR` (dagger ID's)
    - we already have the database client elsewhere
    - we only need the type and the AutoMigrate function (via the singleton Client)


EID: Environment ID and Table

- `eid` - sha256 hash of value content
- `kind` - enum, "directory" | "container"
- `dagger_id` - VARCHAR because they can get long
    