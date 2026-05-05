# Runtime Services

This directory contains the core business logic services for the runtime.

## Services

- **[Environ](environ/AGENTS.md)**: Manages virtualized environments backed by [Dagger](https://dagger.io). Provides a filesystem and execution environment for agents.
- **[Session](session/AGENTS.md)**: Manages agent sessions, event history, and state persistence via GORM.
- **Artifact**: Simple blob storage for files referenced by agents or sessions (Filesystem-backed).

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
