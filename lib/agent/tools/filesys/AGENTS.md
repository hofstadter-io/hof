# Filesystem Tools

This directory contains tools for agents to interact with the filesystem of their attached environment.

## Files

- `filesys.go`: Implements a suite of filesystem tools (`fs_read`, `fs_write`, `fs_edit`, `fs_list`, `fs_glob`, `fs_grep`, `fs_del`).
    - **Integration**: These tools interact with the `environ` service (Dagger container).
    - **State Management**: Modifying tools (`write`, `edit`, `del`) update the `currEnv` state with the new environment URI returned by the `environ` service, ensuring the agent always operates on the latest filesystem state.
