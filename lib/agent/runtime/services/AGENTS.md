# Services

This directory contains the core domain services for the agent runtime.

## Files

- `container.go`: Defines the `Container` interface, which abstracts the execution environment and filesystem (backed by Dagger). It provides a unified API for execution (`Exec`) and filesystem operations (`ReadFile`, `WriteFile`, etc.) used by handlers and agents.
- `db.go`: Placeholder for database client/handle management.

## Subdirectories

- `artifact/`: Artifact storage service (filesystem-backed).
- `environ/`: Environment service implementation using Dagger.
- `session/`: Session management service (DB-backed).
