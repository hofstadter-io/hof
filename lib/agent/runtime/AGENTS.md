# Agent Runtime

The `runtime` package implements the core server and application logic for the agent platform. It manages WebSocket connections, virtualized environments, and agent sessions.

## Core Components

### Runtime (`runtime.go`)
The singleton application state container. Initializes and holds references to all services, database connections, and active clients.

```go
type Runtime struct {
	AppName string

	Ctx context.Context
	mu  sync.Mutex // To protect clients map among other things
	db  *gorm.DB
	e   *echo.Echo

	// services
	A artifact.Service
	S session.Service

	// agentic stuff
	Models  map[string]model.LLM
	Agentic agents.Config

	// clients & comms
	Handlers   map[string]Handler
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}
```

### Client (`client.go`)
Represents a connected user agent (e.g., a VS Code window) via WebSocket. Handles the message loop.

```go
type Client struct {
	User  string
	State map[string]any // should this be persisted, do we even need it with user:... State? (same user on two clients, repo in different locations?)

	// when we have custom agents, or local to a session even? (b/c diff sess diff workdir)
	AgentDefs map[string]agents.Agent

	// this really depends on the workspace / session
	// and should also be merged with (1) user global (2) builtin defaults
	// need a place for selecting which ones show up in the dropdown vs @mention [any]
	Agentic agents.Config

	// we should perhaps store active sessions here
	// various information we'd like to share between agents (multiple vscode status/state)

	// other stuff needs to be persisted
	// 1. agent config (maybe we just store these in the state with user:...)
	// 2. session state/history (already done by SessionService, but needs improvements)

	conn *websocket.Conn

	send chan []byte // Buffered channel for outbound messages

	handleMessage func(*Client, *Message)
}

// readPump() handles incoming messages
// writePump() handles outgoing messages
```

### REST API (`api_environ.go`)
Provides HTTP endpoints for filesystem and environment operations, primarily used by the VS Code extension or other clients for synchronous operations.

- `POST /fs/read`
- `POST /fs/write`
- `POST /fs/stat`
- `POST /fs/list`

## Sub-directories

### [Services](services/AGENTS.md)
Contains the business logic for Environments (Dagger), Sessions (DB), and Artifacts.

### Handlers (`handlers/`)
WebSocket message handlers.
- **`ws/`**: Implementations for specific message types (chat, info, etc.).
