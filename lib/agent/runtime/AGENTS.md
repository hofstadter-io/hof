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
	State map[string]any 

	// when we have custom agents, or local to a session even? (b/c diff sess diff workdir)
	AgentDefs map[string]agents.Agent

	// this really depends on the workspace / session
	// and should also be merged with (1) user global (2) builtin defaults
	Agentic agents.Config

	conn *websocket.Conn

	send chan []byte // Buffered channel for outbound messages

	handleMessage func(*Client, *Message)
}
```

### REST API
Provides HTTP endpoints for filesystem and environment operations, primarily used by the VS Code extension or other clients for synchronous operations. These are implemented in `handlers/api/`.

- `POST /fs/*`: Filesystem operations (read, write, stat, list, etc.) backed by the environment service.

### Sub-directories

### [Services](services/AGENTS.md)
Contains the business logic for Environments (Dagger), Sessions (DB), and Artifacts.

### [Handlers](handlers/AGENTS.md)
Implements WebSocket (`ws/`) and REST API (`api/`) handlers, with shared logic in `common/`.
