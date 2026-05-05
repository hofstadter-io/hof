# WebSocket Handlers

This directory contains the WebSocket message handlers for the runtime. These handlers process incoming JSON payloads from connected clients (e.g., VS Code extension) and perform actions like running agents, managing sessions, or retrieving information.

## Handler Map

Handlers are typically registered in a map (seen in `runtime/runtime.go` or `index.go` here) keyed by a string event name (e.g., "chat", "session.create").

## Handlers

### Chat (`chat.go`)
Handles interaction with agents (LLMs).

- **Event**: `chat`
- **Payload**: `common.ChatPayload` (via `chat.go` or `makeChatUserMessageHandler`)
- **Flow**:
    1.  Unmarshals payload.
    2.  Calls `common.SessionChat` to start the agent runner.
    3.  Streams events back to the client via `c.Mail("chat.event", ...)`.

### Sessions (`sessions.go`)
Manages the lifecycle of user sessions via WebSocket events.

- **Events**:
    - `session.list`: Lists all sessions.
    - `session.get`: Gets details for a specific session.
    - `session.create`: Creates a new session.
    - `session.delete`: Removes a session.
    - `session.clone`: Clones a session.
    - `session.splice`: Splices a session's event history.
    - `session.state.*`: Get/Put/Del specific keys in the session state.

*Note: Many of these are being refactored to use `common/` handlers.*

### Info (`info.go`)
Provides system and configuration information.

- **Events**:
    - `config.info`: Returns current agent configuration.
    - `models.list`: Lists available LLM models.
    - `agents.list`: Lists available agent definitions.

### Index (`index.go`)
Registers all WebSocket handlers into the `aruntime.Runtime.Handlers` map.
