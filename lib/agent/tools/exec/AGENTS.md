# Exec Tool

This directory contains the `exec` tool, which allows agents to run shell commands in their attached environment.

## Files

- `exec.go`: Implements the `Exec` tool. It:
    1. Retrieves the current environment URI from the agent's state (`currEnv`).
    2. Uses the `environ` service to execute the provided script in the container.
    3. Updates the `currEnv` state with the new environment URI (since execution changes state).
    4. Returns the command's exit code, stdout, and stderr.
