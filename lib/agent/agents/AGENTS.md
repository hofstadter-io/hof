# Agent Configuration

This directory handles the configuration and instantiation of agents, primarily using [CUE](https://cuelang.org/).

## Key Components

### CUE Mapping (`cue.go`)
Maps declarative CUE configurations to Go types and constructs the runtime agent objects.

- **`AgenticCUE`**: Loads, builds, and validates the CUE configuration from the directory.
- **`BuildAgent`**: Factory function. Takes the config and creates an `llmagent.New()` instance. It handles:
    - Tool initialization (`buildTools`)
    - MCP toolset setup (`buildMcp`)
    - Model selection
    - Callback registration

## Configuration Structure

```go
type Config struct {
	Models   map[string]Model   `json:"models"`
	Agents   map[string]Agent   `json:"agents"`
	Tools    map[string]Tool    `json:"tools"`
	Toolsets map[string]Toolset `json:"toolsets"`
	Environs map[string]Environ `json:"environs"`

	Embeds   map[string]any    `json:"embeds"`
	EmbedDir string            `json:"embedDir"`
	AgentsMD map[string]string `json:"agentsMD"`

	Templates templates.TemplateMap
}

type Agent struct {
	// proxy to adk fields
	Name        string `json:"name"`
	Model       string `json:"model"`
	Description string `json:"description"`
	Instruction string `json:"instruction"`

	Tools     []string `json:"tools"`
	Toolsets  []string `json:"toolsets"`
	Mcp       []string `json:"mcp"`
	SubAgents []string `json:"subagents"`

	// veg concepts, some of this is more tied to the session, but every session starts with an agent
	AutoLoadWorkdir bool              `json:"autoLoadWorkdir"`   // we need a way to say yay/nay to mounting the local dir, we don't need it for many queries
	Environ         string            `json:"environ,omitempty"` // what is the agent default, none means no container
	AgentsMD        map[string]string `json:""`
}
```
