# Agent Runtime

This repository contains the implementation of the Agent Runtime, a server-based system that manages autonomous agents, their sessions, and their execution environments.

## Architecture

The system is designed around a single server that handles multiple client connections. Each client can manage multiple sessions.

### Core Concepts

- **Runtime**: The root object holding core services (DB, Dagger, Server, Clients).
- **Client**: Represents a connected user (e.g., VS Code Extension) via WebSocket.
- **Session**: A persistent interaction context involving agents, history, and state.
- **Environment**: A virtualized, containerized execution environment (backed by Dagger) where agents run commands and manipulate files.

### Communication

- **WebSocket**: Primary channel for async streaming, chat, and event broadcasting.
- **REST**: Used for synchronous operations and file content retrieval.

## Directory Structure

For detailed documentation, refer to the `AGENTS.md` in each directory.

- [**agents/**](agents/AGENTS.md): Configuration mapping (CUE -> Go) for agents.
- [**cmd/**](cmd/AGENTS.md): CLI entrypoint for the runtime.
- [**extension/**](extension/AGENTS.md): Entrypoint for the VS Code extension background server.
- [**models/**](models/AGENTS.md): LLM model initializers (e.g., Gemini).
- [**runtime/**](runtime/AGENTS.md): Core server logic, WebSocket handlers, and domain services.
    - [**handlers/ws/**](runtime/handlers/ws/AGENTS.md): WebSocket event handlers.
    - [**services/**](runtime/services/AGENTS.md): Domain services (Artifact, Environ, Session).
- [**tools/**](tools/AGENTS.md): Tool implementations available to agents.
    - [**browser/**](tools/browser/AGENTS.md): Browser automation tool.
    - [**filesys/**](tools/filesys/AGENTS.md): Filesystem manipulation tools.
    - [**mcp/**](tools/mcp/AGENTS.md): MCP tool integrations.
