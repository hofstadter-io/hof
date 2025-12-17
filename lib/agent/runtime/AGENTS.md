# Runtime

This directory contains the core logic for the agent runtime server. It manages the lifecycle of the application, services, client connections, and HTTP/WebSocket endpoints.

## Files

- `runtime.go`: Defines the `Runtime` struct, which is the central object holding references to services (Artifact, Session, Environ, DB), LLM models, and connected clients. It handles initialization and server setup.
- `run.go`: Implements the `Run` method to start the Echo HTTP server and the client management loop. It also handles the WebSocket upgrade in `serveWs`.
- `client.go`: Defines the `Client` struct, representing a connected user (e.g., VS Code extension). It manages the WebSocket connection (`readPump`, `writePump`) and message routing.
- `api_environ.go`: Implements HTTP endpoints (`/fs/*`) for filesystem and environment operations, delegating to the `environ` service.
- `session.go`: Defines the `Session` struct within the runtime context, likely for managing active session state and associated Dagger containers.

## Subdirectories

- `handlers/`: WebSocket event handlers.
- `services/`: Core domain services (Artifact, Environ, Session).
