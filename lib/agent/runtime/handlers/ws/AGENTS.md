# WebSocket Handlers

This directory contains the WebSocket event handlers for the agent runtime server. These handlers process messages from clients (like the VS Code extension) and manage the interaction with the agent system.

## Files

- `index.go`: The central registry for WebSocket handlers. `SetupHandlers` maps message types (e.g., `chat`, `session.get`) to their respective handler functions.
- `chat.go`: Handles chat interactions. `chatUserMessage` receives user input, initializes the requested agent, creates a runner, and streams the agent's execution events back to the client.
- `info.go`: Provides informational endpoints. Handlers like `configInfo`, `modelsList`, and `agentsList` return configuration details to the client.
- `sessions.go`: Manages user sessions. Includes handlers for creating, listing, retrieving, and deleting sessions, as well as managing session state (key-value store).

## Key Handlers

- `chat`: Main entry point for chatting with an agent.
- `session.create`: Initializes a new session, optionally setting up an environment.
- `session.get` / `session.list`: Retrieves session data.
- `config.reload`: Triggers a reload of the runtime configuration.
