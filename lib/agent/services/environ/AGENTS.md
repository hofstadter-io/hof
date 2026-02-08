# Environ Service (Dagger)

This directory implements the Environment service using [Dagger](https://dagger.io). It provides a virtualized, persistent, and versioned execution environment for agents. Every state change (execution, file modification) results in a new immutable version of the environment.

## Architecture

The service uses a singleton `localEnviron` that holds the Dagger client and a GORM database connection. Environments are stored as Dagger containers and metadata is tracked in a database `environs` table).

- **Versioning**: Environments are versioned using URIs (e.g., `oci://host:5000/uuid:tag`). Operations that modify the environment create a new container and increment the tag.
- **Persistence**: Container states are pushed to a local registry (e.g., `host.docker.internal:5000`) and metadata is saved to the database.

## Key Files

- **`dagger_client.go`**: Service initialization, singleton management, and basic listing functions.
- **`dagger_database.go`**: Database models (`Environ`), migration, and helper functions.
- **`dagger_env_create.go`**: Logic for creating new environments. Supports building from existing OCI images (`FromUri`) and attaching source filesystems.
- **`dagger_exec.go`**: Implements command execution (`Exec`).
- **`dagger_fs_mutate.go`**: Filesystem modification operations (`WriteFile`, `EditFile`, `Delete`, `Copy`, `Move`).
- **`dagger_fs_query.go`**: Filesystem inspection operations (`Stat`, `ReadFile`, `ReadDirectory`).

## Key Concepts

- **Environ URI**: Unique identifier for a specific state of an environment.
- **Tag Increment**: The mechanism for moving from one state to the next (e.g., `...:0` -> `...:1`).

## Usage Example

```go
// From runtime/services/environ/dagger_client.go
// Initialize and use
err := environ.Initialize(ctx, db)
client := environ.Client()
files, err := client.ReadDirectory(uri, path, diff)
```
