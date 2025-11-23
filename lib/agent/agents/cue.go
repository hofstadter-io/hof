package agents

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/hofstadter-io/hof/lib/agent/tools/filesys"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
)

type Config struct {
	Agents map[string]Agent `json:"agents"`
}

type Agent struct {
	Name        string `json:"name"`
	Model       string `json:"model"`
	Description string `json:"description"`
	Instruction string `json:"instruction"`

	Tools     []string `json:"tools"`
	SubAgents []string `json:"subagents"`
}

// this code constructs one or more agents from a CUE value
// to build up an agentic system
func AgenticCUE(agentDir string, models map[string]model.LLM) ([]agent.Agent, error) {
	// loadup and validate our agentic CUE
	ctx := cuecontext.New()
	entrypoints := []string{agentDir}
	bis := load.Instances(entrypoints, nil)
	bi := bis[0]
	if bi.Err != nil {
		return nil, fmt.Errorf("while loading agentic CUE: %w", bi.Err)
	}
	val := ctx.BuildInstance(bi)
	if val.Err() != nil {
		return nil, fmt.Errorf("while building agentic CUE: %w", val.Err())
	}

	if err := val.Validate(); err != nil {
		return nil, fmt.Errorf("while validating agentic CUE: %w", err)
	}

	fmt.Println("AgenticCUE.value:", val, "\n\n")

	// decode the agentic CUE into a struct
	var config Config
	err := val.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("while decoding agentic CUE: %w", err)
	}
	// fmt.Println("AgenticCUE.config:", config)

	agents := []agent.Agent{}
	for _, a := range config.Agents {
		A, err := buildAgent(config, a.Name, models)
		if err != nil {
			return nil, fmt.Errorf("while building agent %q: %w", a.Name, err)
		}
		agents = append(agents, A)
	}

	return agents, nil
}

func buildAgent(config Config, agentName string, models map[string]model.LLM) (agent.Agent, error) {
	agent := config.Agents[agentName]
	model, ok := models[agent.Model]
	if !ok {
		return nil, fmt.Errorf("unknown model %q in agent %q", agent.Model, agent.Name)
	}

	c := llmagent.Config{
		Name:        agent.Name,
		Model:       model,
		Description: agent.Description,
		Instruction: agent.Instruction,
	}

	for _, t := range agent.Tools {
		var T tool.Tool
		var err error
		switch t {
		case "directory_tree":
			T, err = filesys.NewTreeDir()

		case "list_directory":
			T, err = filesys.NewReadDir()

		case "grep_regexp":
			T, err = filesys.NewGrepFiles()

		case "read_file":
			T, err = filesys.NewReadFile()

		case "write_file":
			T, err = filesys.NewWriteFile()

		default:
			if agentAsTool, found := strings.CutPrefix(t, "@"); found {
				A, aerr := buildAgent(config, agentAsTool, models)
				if aerr != nil {
					return nil, fmt.Errorf("error creating agent tool %q in agent %q", t, agent.Name)
				}
				T = agenttool.New(A, &agenttool.Config{
					SkipSummarization: true,
				})
			} else {
				return nil, fmt.Errorf("unknown tool %q in agent %q", t, agent.Name)
			}
		}

		if err != nil {
			return nil, fmt.Errorf("while creating tool %s: %w", t, err)
		}
		c.Tools = append(c.Tools, T)
	}

	for _, sa := range agent.SubAgents {
		if subagent, found := strings.CutPrefix(sa, "@"); found {
			A, aerr := buildAgent(config, subagent, models)
			if aerr != nil {
				return nil, fmt.Errorf("error creating agent subagent %q in agent %q", subagent, agent.Name)
			}
			c.SubAgents = append(c.SubAgents, A)
		} else {
			return nil, fmt.Errorf("unknown subagent %q in agent %q", sa, agent.Name)
		}
	}

	return llmagent.New(c)
}
