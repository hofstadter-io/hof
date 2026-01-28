package config

import (
	"cuelang.org/go/cue"
	"github.com/hofstadter-io/hof/lib/templates"
)

type Config struct {
	Models   map[string]Model   `json:"models"`
	Agents   map[string]Agent   `json:"agents"`
	Tools    map[string]Tool    `json:"tools"`
	Toolsets map[string]Toolset `json:"toolsets"`
	Environs map[string]Environ `json:"environs"`

	Embeds   map[string]any `json:"embeds"`
	EmbedDir string         `json:"embedDir"`

	// agents instruction files not tied to a project / dir / env?
	AgentsMD map[string]string `json:"agentsMD"`
	// todo, skills and all that

	Templates templates.TemplateMap `json:"-"`
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
	AutoLoadWorkdir bool   `json:"autoLoadWorkdir"`   // we need a way to say yay/nay to mounting the local dir, we don't need it for many queries
	Environ         string `json:"environ,omitempty"` // what is the agent default, none means no container

	// agent specific extra md files, eventually more?
	AgentsMD map[string]string `json:"-"`
}

type Model struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Toolset struct {
	Name  string `json:"name"`
	Tools []Tool `json:"tools"`
}

type AgentMD struct {
	Path     string         `json:"path"`
	Content  string         `json:"content"`
	Metadata map[string]any `json:"metadata"`
}

type Environ struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	Spec      EnvironSpec `json:"spec,omitempty"`
	SpecValue cue.Value   `json:"specValue,omitempty"`
}

type EnvironSpec struct {
	From       string            `json:"from,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	Workdir    string            `json:"workdir,omitempty"`
	Entrypoint []string          `json:"entrypoint,omitempty"`
	Ports      map[string][]int  `json:"ports,omitempty"`
	User       string            `json:"user,omitempty"`
}
