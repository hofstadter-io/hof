package agents

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/genai"

	"github.com/hofstadter-io/hof/lib/agent/tools/meta"
	"github.com/hofstadter-io/hof/lib/templates"
)

type Config struct {
	Models map[string]Model `json:"models"`
	Agents map[string]Agent `json:"agents"`
	Tools  map[string]Tool  `json:"tools"`

	Instructions    map[string]any `json:"instructions"`
	InstructionsDir string         `json:"instructionsDir"`

	Templates templates.TemplateMap
}

type Agent struct {
	Name        string `json:"name"`
	Model       string `json:"model"`
	Description string `json:"description"`
	Instruction string `json:"instruction"`

	Tools     []string `json:"tools"`
	SubAgents []string `json:"subagents"`
}

type Model struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Instruction string `json:"instruction"`
}

// this code constructs one or more agents from a CUE value
// to build up an agentic system
func AgenticCUE(agentDir string, models map[string]model.LLM) (config Config, err error) {
	// loadup and validate our agentic CUE
	ctx := cuecontext.New()
	entrypoints := []string{agentDir}
	bis := load.Instances(entrypoints, nil)
	bi := bis[0]
	if bi.Err != nil {
		return config, fmt.Errorf("while loading agentic CUE: %w", bi.Err)
	}
	val := ctx.BuildInstance(bi)
	if val.Err() != nil {
		return config, fmt.Errorf("while building agentic CUE: %w", val.Err())
	}

	if err := val.Validate(); err != nil {
		return config, fmt.Errorf("while validating agentic CUE: %w", err)
	}

	// fmt.Println("AgenticCUE.value:", val, "\n\n")

	// decode the agentic CUE into a struct
	err = val.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("while decoding agentic CUE: %w", err)
	}

	err = prepareTemplates(&config)
	if err != nil {
		return config, fmt.Errorf("while preparing templates: %w", err)
	}

	// fmt.Println("AgenticCUE.config:", config)
	return config, nil
}

func BuildAgent(config Config, agentName, modelName string, models map[string]model.LLM) (agent.Agent, error) {
	agt := config.Agents[agentName]
	if modelName == "" || modelName == "default" {
		modelName = agt.Model
	}
	fmt.Println("BuildAgent", agentName, modelName)
	mdl, ok := models[modelName]
	if !ok {
		return nil, fmt.Errorf("unknown model %q in agent %q", modelName, agt.Name)
	}

	c := llmagent.Config{
		Name:        agt.Name,
		Model:       mdl,
		Description: agt.Description,
		// Instruction:         agent.Instruction,
		InstructionProvider: renderInstructions(config, agt),
	}

	ts, err := buildTools(config, agt, models)
	if err != nil {
		return nil, fmt.Errorf("while building tools for %q: %w", agt.Name, err)
	}
	c.Tools = ts

	addCallbacks(config, agt, &c)

	for _, sa := range agt.SubAgents {
		if subagent, found := strings.CutPrefix(sa, "@"); found {
			A, aerr := BuildAgent(config, subagent, "default", models)
			if aerr != nil {
				return nil, fmt.Errorf("error creating agent subagent %q in agent %q", subagent, agt.Name)
			}
			c.SubAgents = append(c.SubAgents, A)
		} else {
			return nil, fmt.Errorf("unknown subagent %q in agent %q", sa, agt.Name)
		}
	}

	return llmagent.New(c)
}

func buildTools(cfg Config, agt Agent, models map[string]model.LLM) ([]tool.Tool, error) {
	var ts []tool.Tool
	for _, t := range agt.Tools {
		fmt.Printf("%s.tool: %q\n", agt.Name, t)
		var T tool.Tool
		var err error

		// @<agent> handling
		agentAsTool, found := strings.CutPrefix(t, "@")
		fmt.Printf("%s.tool.agent: %q ? %v\n", agt.Name, agentAsTool, found)
		if found {
			A, aerr := BuildAgent(cfg, agentAsTool, "default", models)
			if aerr != nil {
				return nil, fmt.Errorf("error creating agent tool %q in agent %q: %w", t, agt.Name, aerr)
			}
			T = agenttool.New(A, &agenttool.Config{
				SkipSummarization: true,
			})
			ts = append(ts, T)
			continue
		}

		// otherwise a builtin tool
		tcfg, ok := cfg.Tools[t]
		if !ok {
			return nil, fmt.Errorf("unknown tool %q in agent %q", t, agt.Name)
		}
		switch t {

		case "cache_put", "cache_write":
			T, err = meta.CacheWrite(tcfg.Name, tcfg.Description)
		case "cache_edit":
			T, err = meta.CacheEdit(tcfg.Name, tcfg.Description)
		case "cache_del", "cache_remove":
			T, err = meta.CacheRemove(tcfg.Name, tcfg.Description)

		case "fs_read":
			T, err = meta.FilesysRead(tcfg.Name, tcfg.Description)
		case "fs_list":
			T, err = meta.FilesysList(tcfg.Name, tcfg.Description)
		case "fs_grep":
			T, err = meta.FilesysGrep(tcfg.Name, tcfg.Description)
		case "fs_edit":
			T, err = meta.FilesysEdit(tcfg.Name, tcfg.Description)
		case "fs_write":
			T, err = meta.FilesysWrite(tcfg.Name, tcfg.Description)
		case "fs_del":
			T, err = meta.FilesysDel(tcfg.Name, tcfg.Description)

		default:
			return nil, fmt.Errorf("unknown tool %q in agent %q %q %v", t, agt.Name, agentAsTool, found)
		}
		if err != nil {
			return nil, fmt.Errorf("while creating tool %s: %w", t, err)
		}

		// keep the tool
		ts = append(ts, T)
	}

	fmt.Println("Final Tools:", ts)
	return ts, nil
}

func addCallbacks(config Config, agt Agent, c *llmagent.Config) {
	c.BeforeAgentCallbacks = []agent.BeforeAgentCallback{
		func(ctx agent.CallbackContext) (*genai.Content, error) {
			fmt.Printf("\nBAC.%s\n", ctx.AgentName())
			return nil, nil
		},
	}

	c.BeforeModelCallbacks = []llmagent.BeforeModelCallback{
		func(ctx agent.CallbackContext, req *model.LLMRequest) (*model.LLMResponse, error) {
			fmt.Printf("\nBMC.%s\n", ctx.AgentName())

			// print system prompt before sending to LLM
			fmt.Println(req.Config.SystemInstruction.Parts[0].Text)
			return nil, nil
		},
	}

	c.BeforeToolCallbacks = []llmagent.BeforeToolCallback{
		func(ctx tool.Context, t tool.Tool, args map[string]any) (map[string]any, error) {
			fmt.Printf("\nBTC.%s.%s %v\n", ctx.AgentName(), t.Name(), args)
			return nil, nil
		},
	}

	//
	// reverse order on the way out
	//

	c.AfterToolCallbacks = []llmagent.AfterToolCallback{
		func(ctx tool.Context, t tool.Tool, args, result map[string]any, err error) (map[string]any, error) {
			fmt.Printf("\nATC.%s.%s %v %v\n", ctx.AgentName(), t.Name(), args, result)
			return result, err
		},
	}

	c.AfterModelCallbacks = []llmagent.AfterModelCallback{
		func(ctx agent.CallbackContext, res *model.LLMResponse, err error) (*model.LLMResponse, error) {
			fmt.Printf("\nAMC.%s\n%#+v\n", ctx.AgentName(), res)
			return res, err
		},
	}

	c.AfterAgentCallbacks = []agent.AfterAgentCallback{
		func(ctx agent.CallbackContext) (*genai.Content, error) {
			fmt.Printf("\nAAC.%s\n", ctx.AgentName())
			return nil, nil
		},
	}
}

func prepareTemplates(config *Config) error {

	cwd, _ := os.Getwd()
	// todo, also put this on the Session
	dir := filepath.Join(cwd, config.InstructionsDir)
	pre := strings.TrimSuffix(dir, "/**/*.md")
	config.Templates = templates.NewTemplateMap()
	// fmt.Printf("found %d templates in %q %q\n", len(config.Templates), dir, pre)
	err := config.Templates.ImportFromFolder(dir, pre, templates.Delims{}, nil)
	if err != nil {
		return fmt.Errorf("while loading instruction templates (%s,%s): %w", cwd, config.InstructionsDir, err)
	}
	fmt.Printf("found %d templates in %s\n", len(config.Templates), dir)

	for _, T1 := range config.Templates {
		for _, T2 := range config.Templates {
			if T1.Name == T2.Name {
				continue
			}
			t := T1.T.New(T2.Name)
			_, err := t.Parse(T2.Source)
			if err != nil {
				return fmt.Errorf("while cross registering templates (%s,%s): %w", T1.Name, T2.Name, err)
			}
		}

		// fmt.Println(T1.Name)
		// for _, t := range T1.T.Templates() {
		// 	fmt.Printf(" - %s\n", t.Name())
		// }
	}

	return nil
}

func renderInstructions(cfg Config, agt Agent) llmagent.InstructionProvider {

	return func(ctx agent.ReadonlyContext) (string, error) {
		// TODO, this last arg is annoying, should have two funcs
		fmt.Println("renderInstructions.Agent", agt.Name)

		var err error
		var t *templates.Template

		t, ok := cfg.Templates[agt.Instruction]
		if !ok {
			// load instruction template
			t, err = templates.CreateFromString(agt.Name, agt.Instruction, templates.Delims{})
			if err != nil {
				fmt.Println("ERROR.renderInstructions.Create", err)
				return "", err
			}
			t.Name = agt.Name + "-inline"
			for _, T := range cfg.Templates {
				t := t.T.New(T.Name)
				_, err := t.Parse(T.Source)
				if err != nil {
					return "", fmt.Errorf("while cross registering templates (%s,%s): %w", t.Name(), T.Name, err)
				}
			}
		}

		// gather data
		data, err := prepareData(cfg, agt)(ctx)
		if err != nil {
			fmt.Println("ERROR.renderInstructions.Prepare", err)
			return "", err
		}

		// render instruction (first time) to get length
		b, err := t.Render(data)
		if err != nil {
			fmt.Println("ERROR.renderInstructions.Render.First", err)
			return "", err
		}

		if strings.Contains(string(b), "CONTEXT SIZE:") {
			data["contextSize"] = len(b)

			b, err = t.Render(data)
			if err != nil {
				fmt.Println("ERROR.renderInstructions.Render.Final", err)
				return "", err
			}
		}

		s := string(b)

		// TODO, add conditional logging from Agent config
		// fmt.Printf("renderInstructions.Final %s\n%s\n", agt.Name, s)

		return s, nil
	}
}

func prepareData(cfg Config, agt Agent) func(ctx agent.ReadonlyContext) (map[string]any, error) {
	return func(ctx agent.ReadonlyContext) (map[string]any, error) {
		data := make(map[string]any)

		// environment of the workspace / vscode
		state := maps.Collect(ctx.ReadonlyState().All())
		data["env"] = map[string]any{
			"basedir": state["basedir"],
		}
		data["config"] = cfg
		data["agent"] = agt

		// agent cache
		cache := make(map[string]any)
		for k, v := range state {
			if p, matched := strings.CutPrefix(k, fmt.Sprintf("cache:%s:", ctx.AgentName())); matched {
				switch p {
				case "planning":
					data["planning"] = v
				case "subconscious":
					data["subconscious"] = v

				default:
					cache[p] = v
				}
			}
		}
		data["cache"] = cache

		stateKeys := slices.Collect(maps.Keys(state))
		cacheKeys := slices.Collect(maps.Keys(cache))
		dataKeys := slices.Collect(maps.Keys(data))

		fmt.Println("stateKeys:")
		for _, k := range stateKeys {
			fmt.Println(" ", k)
		}
		fmt.Println("cacheKeys:")
		for _, k := range cacheKeys {
			fmt.Println(" ", k)
		}
		fmt.Println("dataKeys:")
		for _, k := range dataKeys {
			fmt.Println(" ", k)
		}
		fmt.Println("subconscious:", data["subconscious"])

		// b, err := json.MarshalIndent(data["cache"], "", "  ")
		// if err != nil {
		// 	fmt.Println("error while marshalling data for debug of template input:", err)
		// }
		// fmt.Println(string(b))

		return data, nil
	}
}
