# Handlers

This directory contains event handlers for the runtime server, supporting both WebSocket and REST API interfaces.

## Common Logic

This layer provides shared business logic and session management used by both `ws/` and `api/` handlers to ensure consistency and avoid duplication.

### Chat (`common/chat.go`)
Implements the core agent interaction loop.

```go
func SessionChat(r *runtime.Runtime, ar *aruntime.Runtime, p *ChatPayload) (*aruntime.Session, error)
```

**Payload**:
```go
type ChatPayload struct {
	User    string `json:"user"`
	Text    string `json:"text"`
	Sid     string `json:"sid"`
	Agent   string `json:"agent"`
	Model   string `json:"model"`
	Environ string `json:"environ"`
}
```

### Session Lifecycle (`common/create.go`, `common/session.go`)
Functions for creating, retrieving, listing, and deleting sessions.

```go
func SessionCreate(r *runtime.Runtime, ar *aruntime.Runtime, payload CreatePayload) (session.Session, error)
func SessionGet(r *runtime.Runtime, ar *aruntime.Runtime, sid string) (session.Session, error)
func SessionList(r *runtime.Runtime, ar *aruntime.Runtime) ([]session.Session, error)
func SessionDel(r *runtime.Runtime, ar *aruntime.Runtime, sid string) error
```

### Session Operations (`common/session_ops.go`)
Advanced session manipulations like cloning, splicing, and state management.

```go
func SessionClone(ctx context.Context, ar *aruntime.Runtime, sid string, pos int) (session.Session, error)
func SessionSplice(ctx context.Context, ar *aruntime.Runtime, sid string, pos, count int, fill []*session.Event) (session.Session, error)
func SessionStateGet(ctx context.Context, ar *aruntime.Runtime, sid, key string) (any, error)
func SessionStatePut(ctx context.Context, ar *aruntime.Runtime, sid, key string, val any) error
func SessionStateDel(ctx context.Context, ar *aruntime.Runtime, sid, key string) error
func SessionPromptRender(ctx context.Context, ar *aruntime.Runtime, sid, agentName string) (string, error)
```

### Filesystem (`common/filesys.go`)
Unified filesystem logic used by both REST and WebSocket handlers. Automatically handles session state updates and access verification.
System-level operations like reloading configuration and listing available environments.

```go
func ReloadConfig(ar *aruntime.Runtime) error
func ListEnvirons() ([]string, error)
```

## Sub-directories

### [REST API (`api/`)](./api/AGENTS.md)
Standard HTTP endpoints for synchronous operations, such as filesystem access, environment listing, and session manipulations.

### [WebSocket (`ws/`)](./ws/AGENTS.md)
Handlers for real-time communication via WebSockets, primarily used by the VS Code extension for chat and live state updates.

## Core Concepts

- **Unified Logic**: Most session and agent operations are being migrated to the `common/` package.
- **Protocol Independence**: The `common/` layer is agnostic to whether the request came via WebSocket or HTTP.
- **Runtime Integration**: Handlers interact with the central `aruntime.Runtime` and `session.Service`.
- **URI Conventions**: 
    - **`veg://`**: Used within VS Code to identify files in the virtual filesystem. The authority and path segments typically encode the session/environment ID and version.
    - **`oci://`**: Used by the backend services (Dagger) to uniquely identify specific versions of an environment. 
    - Translation occurs in the VS Code extension (`vsUriToVeg`) before making REST API calls.
