# Extension Entrypoint

This directory contains the entrypoint for the VS Code extension's background server.

## Files

- `cmd.go`: Provides the `Run` function specifically for the VS Code extension context. It initializes the runtime and starts the WebSocket server, acting as the bridge between the extension and the Go agent runtime.
