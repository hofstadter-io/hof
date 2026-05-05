# Tools

This directory contains the tools available to agents. Tools are split into core capabilities and specialized integrations.

## Core Capabilities

### Filesystem (`filesys/`)
File manipulation tools interacting with the agent's environment (`environ.Client()`).

- **Tools**: `fs_read`, `fs_list`, `fs_glob`, `fs_grep`, `fs_write`, `fs_edit`, `fs_del` (Names derived from implementation, exposed as function tools).
- **Implementation**: `filesys/filesys.go`

```go
// Example: FilesysRead
func FilesysRead(name, description string) (tool.Tool, error) {
    handler := func(ctx tool.Context, input FilesysPathArgs) (FilesysResult, error) {
        // ... gets content from environ.Client().ReadFile ...
        return FilesysResult{Path: input.Path, Status: "ok"}, nil
    }
    // ...
}

type FilesysResult struct {
	Path   string `json:"path"`            // filesystem path
	Status string `json:"status"`          // "ok" or "error"
	Error  string `json:"error,omitempty"` // the error message if present
}
```

### Execution (`exec/`)
Runs shell scripts or commands in the attached environment.

- **Tools**: `exec`
- **Implementation**: `exec/exec.go`

```go
type ExecArgs struct {
	Script string `json:"script"` // command or script to run
}
type ExecResult struct {
	ExitCode int    `json:"exitCode"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Status   string `json:"status"`          // "ok" or "error"
	Error    string `json:"error,omitempty"` // error message if there is an error
}
```

### Cache (`cache/`)
Manages the agent's persistent key/value state.

- **Tools**: `cache_write`, `cache_remove`, `cache_edit`
- **Implementation**: `cache/cache.go`

```go
type CacheEditArgs struct {
	Key string `json:"key"` // path to a file
	Old string `json:"old_string"`
	New string `json:"new_string"`
	Exp int    `json:"expected replacements"` // defaults to 1 if not set
}
```

### Mathy (`mathy.go`)
Simple arithmetic tools for demonstration/testing.
- **Tools**: `sum`, `sub`

## Specialized Integrations

### [Browser](browser/AGENTS.md)
TypeScript-based Playwright browser automation.

### [MCP](mcp/AGENTS.md)
Integrations with the Model Context Protocol (Github, Tavily, Local).
