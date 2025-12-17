# Runtime Services

This directory contains the core business logic services for the runtime.

## Services

### Environ (`environ/`)
Manages virtualized environments backed by [Dagger](https://dagger.io).
Provides a filesystem and execution environment for agents.

- **Key File**: `dagger_client.go` - Singleton client initialization.
- **Usage**: `environ.Client()` returns the global `localEnviron` instance.

```go
// Example: Initialize and use
err := environ.Initialize(ctx, db)
client := environ.Client()
files, err := client.ReadDirectory(uri, path, diff)
```

### Session (`session/`)
Manages agent sessions, event history, and state persistence.
Implements `google.golang.org/adk/session.Service`.

- **Implementation**: `databaseService` (GORM-backed).
- **Key Features**:
  - `Create`, `Get`, `List`, `Delete` sessions.
  - `AppendEvent`: Adds events to history and updates state.
  - State management: Merges App, User, and Session state.

```go
// Example: Creating a session
s, err := session.NewSessionServiceGorm(db)
resp, err := s.Create(ctx, &session.CreateRequest{
    AppName: "veg",
    UserID:  "user-1",
})
```

### Artifact (`artifact/`)
Simple blob storage for files referenced by agents or sessions.
- **Implementation**: Filesystem-backed service.

## Interfaces

### Container (`container.go`)
Abstracts the interaction with an environment (Dagger or otherwise).
Defines methods for filesystem operations and execution.

```go
type Container interface {
	// Dagger refs
	ID() (string, error) // hash ref to a dagger ID stored else where
	Load(id string) (*dagger.Container, error)

	// Exec related
	Exec(args ...string) (ExecResult, error)

	// VS Code
	Stat(path string) (FileStat, error)
	ReadFile(path string) (string, error)
	WriteFile(path, content string) error

	CreateDirectory(path string) error
	ReadDirectory(path string) ([]Dirent, error)

	Rename(src, dst string) error
	Copy(src, dst string) error
	Delete(path string, recursive bool) error

	Watch(path string, recursive bool)
}
```
