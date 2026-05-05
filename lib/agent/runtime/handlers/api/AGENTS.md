# REST API Handlers

This directory contains the REST API handlers for the agent runtime, implemented using the Echo framework.

## Core Components

### Runtime (`runtime.go`)
The `Runtime` struct holds the application state. Routes are registered in the `Setup` function.

```go
type Runtime struct {
	AppName string
	S       session.Service
	Agentic *config.Config
}
```

## Handlers

### Filesystem (`filesys.go`)
Handles synchronous filesystem operations (open, stat, read, list, write, delete, etc.) on the virtualized environment. These handlers use `common` handlers to ensure session state and access control.

### Environment (`env.go`)
Lists available environments.

### Sessions (`session.go`)
Handles session-specific operations like cloning, splicing, and instruction rendering.
*Note: These are being refactored to use `common/` handlers.*
