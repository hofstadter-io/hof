# Agent Configuration

This directory contains logic for mapping CUE configurations to Go types for the agent system.

## Files

- `cue.go`: Defines the configuration structures (`Config`, `Agent`, `Tool`, `Environ`) and provides the `AgenticCUE` function to load and validate agent configurations from CUE files. It acts as the bridge between the CUE definition language and the internal Go representation of agents.

## Key Types

- `Config`: Root configuration structure containing agents, models, tools, etc.
- `Agent`: Defines an agent's properties, including tools, sub-agents, and LLM model.
- `Environ`: Defines the execution environment (e.g., Docker container specs) for an agent.
