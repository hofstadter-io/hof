# Cache Tools

This directory contains tools for agents to interact with their own state/cache.

## Files

- `cache.go`: Implements the `cache_put` (write), `cache_del` (remove), and `cache_edit` tools. These allow agents to store and retrieve data in their persistent state context, which is useful for maintaining memory across turns.
