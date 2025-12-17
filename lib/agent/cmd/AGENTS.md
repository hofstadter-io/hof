# Command Line Entrypoint

This directory contains the main CLI entrypoint for running the agent runtime.

## Files

- `run.go`: Provides the `Run` function which bootstraps the agent system. It:
    1. Initializes the `runtime`.
    2. Sets up WebSocket handlers.
    3. Starts the runtime execution loop.
