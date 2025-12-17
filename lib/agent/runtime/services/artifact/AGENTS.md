# Artifact Service

This directory contains implementations for artifact storage, allowing agents to persist and retrieve files or data blobs.

## Files

- `filesys.go`: Implements the `Artifacts` interface using the local filesystem. It provides methods to read, write, and delete artifacts stored in a specific root directory (e.g., `.agents/artifacts`).
