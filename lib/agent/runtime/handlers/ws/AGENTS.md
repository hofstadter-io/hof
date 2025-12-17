# WebSocket Handlers

This directory contains the WebSocket message handlers for the runtime. These handlers process incoming JSON payloads from connected clients (e.g., VS Code extension) and perform actions like running agents, managing sessions, or retrieving information.

## Handler Map

Handlers are typically registered in a map (seen in `runtime/runtime.go` or `index.go` here) keyed by a string event name (e.g., "chat", "session.create").

## Handlers

### Chat (`chat.go`)
Handles interaction with agents (LLMs).

- **Event**: `chat`
- **Payload**: `ChatPayload`
- **Flow**:
    1.  Unmarshals payload (Text, Agent, Model).
    2.  Retrieves the Session.
    3.  Resolves the Environment and Agent Instructions.
    4.  Builds the Agent (`agents.BuildAgent`).
    5.  Runs the Agent using `runner.New()` and `R.Run()`.
    6.  Streams events back to the client via `c.Mail("chat.event", ...)`.

```go
type ChatPayload struct {
	Text  string `json:"text"`
	Sid   string `json:"sid"`
	Agent string `json:"agent"`
	Model string `json:"model"`
}
```

### Sessions (`sessions.go`)
Manages the lifecycle of user sessions.

- **Events**:
    - `session.list`: Lists all sessions for the user.
    - `session.get`: Gets full state/history for a specific session.
    - `session.create`: Creates a new session, optionally with a Dagger environment.
    - `session.delete`: Removes a session.
    - `session.state.*`: Get/Put/Del specific keys in the session state.

```go
type SessionCreateRequest struct {
	Title   string                        `json:"title,omitempty"`
	Focus   bool                          `json:"focus,omitempty"`
	Environ *environ.EnvironCreateOptions `json:"environ,omitempty"`
}
```

### Info (`info.go`)
*Documentation inferred from filename* - Likely handles general system info or status checks.

### Index (`index.go`)
*Documentation inferred from filename* - Likely contains the registration logic mapping string keys to these handler functions.
