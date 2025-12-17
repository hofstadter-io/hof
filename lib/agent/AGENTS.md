# Veggie Codebase Guide

This repository uses a recursive documentation strategy. Each significant directory contains an `AGENTS.md` file explaining its contents, key types, and workflows.

**🤖 Agent Tip**: Always prefer reading the local `AGENTS.md` when entering a directory to understand the context before exploring code files.

## High-Level Architecture

The system is a backend for an AI coding agent (VS Code extension).

1.  **Entrypoint**: `extension/cmd.go` starts the server.
2.  **Runtime**: `runtime/` manages the application state, WebSocket connections (`Client`), and core services.
3.  **Services**: `runtime/services/` handles:
    *   **Environments**: Virtualized filesystems/terminals via [Dagger](https://dagger.io).
    *   **Sessions**: User state and history via SQLite/Postgres.
4.  **Configuration**: `agents/` loads agent definitions, tools, and models using [CUE](https://cuelang.org/).
5.  **Tools**: `tools/` provides capabilities like File I/O, Shell Execution, Browser Automation, and MCP integrations.

## Directory Index

### Core Logic
- **[`runtime/`](runtime/AGENTS.md)**: The heart of the application.
    - **[`runtime/services/`](runtime/services/AGENTS.md)**: Business logic (Dagger, DB).
- **[`agents/`](agents/AGENTS.md)**: Agent configuration and loading logic.

### Capabilities
- **[`tools/`](tools/AGENTS.md)**: Tool implementations.
    - **[`tools/browser/`](tools/browser/AGENTS.md)**: Playwright browser.
    - **[`tools/mcp/`](tools/mcp/AGENTS.md)**: MCP integrations (Github, Tavily).

### Entrypoints & Misc
- **`extension/`**: Contains the main `cmd.go` entrypoint for the VS Code extension.
- **`cmd/`**: CLI entrypoints (e.g., `run.go` for standalone runs).
- **`models/`**: LLM client wrappers (e.g., `gemini.go`).

## Key Concepts

- **ADK (Agent Development Kit)**: The core middle piece that manages sessions, agent/llm requests, and the event system. It is the central library (`google.golang.org/adk`) used throughout the runtime.
- **Dagger**: Used to virtualize the environment the agent interacts with.
- **MCP (Model Context Protocol)**: Standard interface for connecting tools.
- **CUE**: Declarative language for configuring agent behaviors and environments.
