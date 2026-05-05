package agents

import (
	"context"
	"fmt"
	"maps"
	"path"
	"slices"
	"sort"
	"strings"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/genai"

	"github.com/hofstadter-io/hof/lib/agent/config"
	"github.com/hofstadter-io/hof/lib/agent/tools/cache"
	"github.com/hofstadter-io/hof/lib/agent/tools/exec"
	"github.com/hofstadter-io/hof/lib/agent/tools/filesys"
	"github.com/hofstadter-io/hof/lib/agent/tools/mcp"
	"github.com/hofstadter-io/hof/lib/templates"
)

func LoadAgent(cfg *config.Config, agentName string) (config.Agent, error) {
	agt, ok := cfg.Agents[agentName]
	if !ok {
		return config.Agent{}, fmt.Errorf("agent not found: %q", agentName)
	}

	// merge MDs: global < agent
	mds := make(map[string]string)
	maps.Copy(mds, cfg.AgentsMD) // from config (global)
	maps.Copy(mds, agt.AgentsMD) // from agent  (per-agent)
	agt.AgentsMD = mds

	return agt, nil
}

// this needs to be supported through a hierarchy of unification
// dir, project, user, org... with modules and per-request
// (hence the CUE, still todo for more CUEism in memory ^^)

// just cause this file is open... random thought
//
// 1. I have left stuff like this all over the code, should build a specialized agent for this
// 2. Why not build an agent (team) that can...
//   1. search a dir or repo for them (3.1 i.e.), summarize
//   2. do some deep research
//   3. build a plan to tackle them, output something structured
//   4. update roadmap / kanban
// 3. sub-team / agent
//   1. process one file at a time
//   2. store comment and context
//   3. give back to 2.1
//
// So then, can we build an agent that can assemble different setups like this, depending on the task?

// this code constructs one or more agents from a CUE value
// to build up an agentic system
func OldAgenticCUE(agentDir string, models map[string]model.LLM) (cfg *config.Config, err error) {
	// loadup and validate our agentic CUE
	if strings.HasPrefix(agentDir, ".veg") {
		agentDir = "./" + agentDir
	}
	// fmt.Println("AgenticCUE", agentDir)
	ctx := cuecontext.New()
	entrypoints := []string{agentDir}
	bis := load.Instances(entrypoints, &load.Config{
		Package: "veg",
	})
	bi := bis[0]
	if bi.Err != nil {
		return cfg, fmt.Errorf("while loading agentic CUE: %w", bi.Err)
	}
	val := ctx.BuildInstance(bi)
	if val.Err() != nil {
		return cfg, fmt.Errorf("while building agentic CUE: %w", val.Err())
	}

	// if err := val.Validate(); err != nil {
	// 	return config, fmt.Errorf("while validating agentic CUE: %w", err)
	// }

	// fmt.Println("AgenticCUE.value:", val, "\n\n")

	// decode the agentic CUE into a struct
	cfg = new(config.Config)
	err = val.Decode(cfg)
	if err != nil {
		return cfg, fmt.Errorf("while decoding agentic CUE: %w", err)
	}

	// err = prepareTemplates(cfg)
	// if err != nil {
	// 	return cfg, fmt.Errorf("while preparing templates: %w", err)
	// }

	if cfg.AgentsMD == nil {
		cfg.AgentsMD = make(map[string]string)
	}

	// fmt.Println("AgenticCUE.config:", config)
	return cfg, nil
}

func BuildAgent(
	cfg *config.Config,
	agentName string,
	modelName string,
	models map[string]model.LLM,
	environMDs map[string]string, // todo, expand the scope of what environData gets passed, could depend on some of the other params
) (agent.Agent, error) {
	// look up agent and set some defaults
	agt, err := LoadAgent(cfg, agentName)
	if err != nil {
		return nil, err
	}

	if modelName == "" || modelName == "default" {
		modelName = agt.Model
	}
	// fmt.Println("BuildAgent", agentName, modelName)
	mdl, ok := models[modelName]
	if !ok {
		return nil, fmt.Errorf("unknown model %q in agent %q", modelName, agt.Name)
	}

	// todo, also handle environment changes?

	c := llmagent.Config{
		Name:        agt.Name,
		Model:       mdl,
		Description: agt.Description,
		// Instruction:         agent.Instruction,
		InstructionProvider: RenderInstructions(cfg, agt, environMDs),
	}

	ts, err := buildTools(cfg, agt, models, environMDs)
	if err != nil {
		return nil, fmt.Errorf("while building tools for %q: %w", agt.Name, err)
	}
	c.Tools = ts

	mcp, err := buildMcp(cfg, agt, models)
	if err != nil {
		return nil, fmt.Errorf("while building mcp toolsets for %q: %w", agt.Name, err)
	}
	c.Toolsets = append(c.Toolsets, mcp...)

	addCallbacks(cfg, agt, environMDs, &c)

	for _, sa := range agt.SubAgents {
		if subagent, found := strings.CutPrefix(sa, "@"); found {
			A, aerr := BuildAgent(cfg, subagent, "default", models, environMDs)
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

func buildMcp(cfg *config.Config, agt config.Agent, models map[string]model.LLM) ([]tool.Toolset, error) {
	var ts []tool.Toolset
	for _, name := range agt.Mcp {
		var (
			t   tool.Toolset
			err error
		)
		switch name {
		case "github":
			t, err = mcp.GithubMCPToolset(context.Background())
		case "tavily":
			t, err = mcp.TavilyMCPToolset(context.Background())
		case "quickbooks":
			t, err = mcp.QuickbooksMCPToolset(context.Background())
		default:
			err = fmt.Errorf("unknown mcp toolset")
		}
		if err != nil {
			return nil, fmt.Errorf("while initializing %s mcp toolset: %w", name, err)
		}
		ts = append(ts, t)
	}
	return ts, nil
}

func buildTools(cfg *config.Config, agt config.Agent, models map[string]model.LLM, environMDs map[string]string) ([]tool.Tool, error) {
	var ts []tool.Tool
	for _, t := range agt.Tools {
		// fmt.Printf("%s.tool: %q\n", agt.Name, t)
		var T tool.Tool
		var err error

		// @<agent> handling
		agentAsTool, found := strings.CutPrefix(t, "@")
		// fmt.Printf("%s.tool.agent: %q ? %v\n", agt.Name, agentAsTool, found)
		if found {
			A, aerr := BuildAgent(cfg, agentAsTool, "default", models, environMDs)
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

		// cache ops
		case "cache_put", "cache_write":
			T, err = cache.CacheWrite(tcfg.Name, tcfg.Description)
		case "cache_edit":
			T, err = cache.CacheEdit(tcfg.Name, tcfg.Description)
		case "cache_del", "cache_remove":
			T, err = cache.CacheRemove(tcfg.Name, tcfg.Description)

		// fs query
		case "fs_read":
			T, err = filesys.FilesysRead(tcfg.Name, tcfg.Description)
		case "fs_list":
			T, err = filesys.FilesysList(tcfg.Name, tcfg.Description)
		case "fs_glob":
			T, err = filesys.FilesysGlob(tcfg.Name, tcfg.Description)
		case "fs_grep":
			T, err = filesys.FilesysGrep(tcfg.Name, tcfg.Description)

		// fs mutate
		case "fs_edit":
			T, err = filesys.FilesysEdit(tcfg.Name, tcfg.Description)
		case "fs_write":
			T, err = filesys.FilesysWrite(tcfg.Name, tcfg.Description)
		case "fs_del":
			T, err = filesys.FilesysDel(tcfg.Name, tcfg.Description)

		// doThings
		case "exec":
			T, err = exec.Exec(tcfg.Name, tcfg.Description)
		// browser
		// search like
		// veg/flow

		default:
			return nil, fmt.Errorf("unknown tool %q in agent %q %q %v", t, agt.Name, agentAsTool, found)
		}
		if err != nil {
			return nil, fmt.Errorf("while creating tool %s: %w", t, err)
		}

		// keep the tool
		ts = append(ts, T)
	}

	// fmt.Println("Final Tools:", ts)
	return ts, nil
}

func addCallbacks(cfg *config.Config, agt config.Agent, environMDs map[string]string, c *llmagent.Config) {
	c.OnModelErrorCallbacks = []llmagent.OnModelErrorCallback{
		func(ctx agent.CallbackContext, llmRequest *model.LLMRequest, llmResponseError error) (*model.LLMResponse, error) {
			// if llmResponseError != nil {
			// 	fmt.Printf("ModelError.%s: %s\n", ctx.AgentName(), llmResponseError)
			// }
			return nil, nil
		},
	}

	c.BeforeAgentCallbacks = []agent.BeforeAgentCallback{
		func(ctx agent.CallbackContext) (*genai.Content, error) {
			// fmt.Printf("\nBAC.%s\n", ctx.AgentName())
			return nil, nil
		},
	}

	c.BeforeModelCallbacks = []llmagent.BeforeModelCallback{
		func(ctx agent.CallbackContext, req *model.LLMRequest) (*model.LLMResponse, error) {
			// fmt.Printf("\nBMC.%s\n", ctx.AgentName())

			// This next section is all about making sure the state is in a good place
			// to match the data we are about to render instructions with
			data, _ := prepareData(cfg, agt, environMDs)(ctx)

			var pfs []string
			if agtmd, ok := data["agentsMd"].([]config.AgentMD); ok {
				for _, am := range agtmd {
					pfs = append(pfs, am.Path)
				}
			}
			sort.Strings(pfs)

			// Update prompt keys in state (only changed ones)
			prefix := "agentmd:" + ctx.AgentName() + ":"
			desiredKeys := make(map[string]bool)
			for _, pf := range pfs {
				desiredKeys[prefix+pf] = true
			}

			// Loop 1: Remove stale agentmd keys for this agent
			for k := range maps.Collect(ctx.State().All()) {
				if strings.HasPrefix(k, prefix) {
					if !desiredKeys[k] {
						ctx.State().Set(k, nil)
					}
				}
			}

			// Loop 2: Add missing agentmd keys for this agent
			for k := range desiredKeys {
				if val, _ := ctx.State().Get(k); val == nil {
					ctx.State().Set(k, "included")
				}
			}

			// print system prompt before sending to LLM

			// hmmm, little utils like this could get spread throughout the code
			// TODO, make a schema somewhere for the various config (cli, system, per-user, per-session, state)
			// showStr, err := ctx.State().Get("showSystemPrompt")
			// fmt.Println("showSystemPrompt.1?", showStr, err)
			// if showStr != nil {
			// 	show, err := strconv.ParseBool(showStr.(string))
			// 	fmt.Println("showSystemPrompt.2?", show, err)
			// 	if err == nil && show {
			// 		fmt.Println(req.Config.SystemInstruction.Parts[0].Text)
			// 	}
			// }

			return nil, nil
		},
	}

	c.BeforeToolCallbacks = []llmagent.BeforeToolCallback{
		func(ctx tool.Context, t tool.Tool, args map[string]any) (map[string]any, error) {
			// fmt.Printf("\nBTC.%s.%s\n", ctx.AgentName(), t.Name())
			return nil, nil
		},
	}

	//
	// reverse order on the way out
	//

	c.AfterToolCallbacks = []llmagent.AfterToolCallback{
		func(ctx tool.Context, t tool.Tool, args, result map[string]any, err error) (map[string]any, error) {
			// fmt.Printf("\nATC.%s.%s %v\n", ctx.AgentName(), t.Name(), err)
			return result, err
		},
	}

	c.AfterModelCallbacks = []llmagent.AfterModelCallback{
		func(ctx agent.CallbackContext, res *model.LLMResponse, err error) (*model.LLMResponse, error) {
			// fmt.Printf("\nAMC.%s %v\n", ctx.AgentName(), err)
			return res, err
		},
	}

	c.AfterAgentCallbacks = []agent.AfterAgentCallback{
		func(ctx agent.CallbackContext) (*genai.Content, error) {
			// fmt.Printf("\nAAC.%s\n", ctx.AgentName())
			return nil, nil
		},
	}
}

func RenderInstructions(cfg *config.Config, agt config.Agent, environMDs map[string]string) llmagent.InstructionProvider {
	return func(ctx agent.ReadonlyContext) (string, error) {
		return RenderInstructionsWithNameAndState(cfg, agt, ctx.AgentName(), maps.Collect(ctx.ReadonlyState().All()), environMDs)
	}
}

func RenderInstructionsWithNameAndState(cfg *config.Config, agt config.Agent, name string, state map[string]any, environMDs map[string]string) (string, error) {
	// TODO, this last arg is annoying, should have two funcs
	// fmt.Println("RenderInstructions.Agent", agt.Name)

	var err error
	var t *templates.Template

	t, ok := cfg.Templates[agt.Instruction]
	if !ok {
		// load instruction template
		t, err = templates.CreateFromString(agt.Name, agt.Instruction, templates.Delims{})
		if err != nil {
			fmt.Println("ERROR.RenderInstructions.Create", err)
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
	data, err := PrepareDataWithNameAndState(cfg, agt, name, state, environMDs)
	if err != nil {
		fmt.Println("ERROR.RenderInstructions.Prepare", err)
		return "", err
	}

	debugPrintData(data)

	// render instruction (first time) to get length
	b, err := t.Render(data)
	if err != nil {
		fmt.Println("ERROR.RenderInstructions.Render.First", err)
		return "", err
	}

	if strings.Contains(string(b), "CONTEXT SIZE:") {
		// TODO, calculate tokens instead?
		data["contextSize"] = len(b)

		b, err = t.Render(data)
		if err != nil {
			fmt.Println("ERROR.RenderInstructions.Render.Final", err)
			return "", err
		}
	}

	s := string(b)

	// TODO, add conditional logging from Agent config
	// fmt.Printf("RenderInstructions.Final %s\n%s\n", agt.Name, s)

	return s, nil
}

type KVPair struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func prepareData(cfg *config.Config, agt config.Agent, environMDs map[string]string) func(ctx agent.ReadonlyContext) (map[string]any, error) {
	return func(ctx agent.ReadonlyContext) (map[string]any, error) {
		return PrepareDataWithNameAndState(cfg, agt, ctx.AgentName(), maps.Collect(ctx.ReadonlyState().All()), environMDs)
	}
}

func PrepareDataWithNameAndState(cfg *config.Config, agt config.Agent, agentName string, state map[string]any, environMDs map[string]string) (map[string]any, error) {
	data := make(map[string]any)

	// environment of the workspace / vscode
	// TODO, this should be included with the expanded environ context (just environMDs for now)
	data["env"] = map[string]any{
		"basedir": state["basedir"],
	}
	data["config"] = cfg
	data["agent"] = agt

	// extract stuff from state
	files := make(map[string]any)
	cache := make(map[string]any)
	for k, v := range state {

		//
		// HMMM, filtering by agent name first probably breaks things, weh should be agnostic to this
		//       so we can switch or pass state between agents
		//

		// files
		if p, matched := strings.CutPrefix(k, fmt.Sprintf("files:%s:", agentName)); matched {
			files[p] = v
			continue
		}

		// cache entries
		if p, matched := strings.CutPrefix(k, fmt.Sprintf("cache:%s:", agentName)); matched {
			switch p {
			case "planning":
				data["planning"] = v
			case "subconscious":
				data["subconscious"] = v

			default:
				cache[p] = v
			}
			continue
		}

	}

	// Sort files by path
	filesSorted := make([]KVPair, 0, len(files))
	for k, v := range files {
		filesSorted = append(filesSorted, KVPair{Key: k, Value: v})
	}
	sort.Slice(filesSorted, func(i, j int) bool {
		return filesSorted[i].Key < filesSorted[j].Key
	})
	data["files"] = filesSorted

	// Sort cache by key
	cacheSorted := make([]KVPair, 0, len(cache))
	for k, v := range cache {
		cacheSorted = append(cacheSorted, KVPair{Key: k, Value: v})
	}
	sort.Slice(cacheSorted, func(i, j int) bool {
		return cacheSorted[i].Key < cacheSorted[j].Key
	})
	data["cache"] = cacheSorted

	agtmd := make(map[string]string)
	for _, f := range filesSorted {
		fpath := f.Key

		// check environMDs (Project files)
		for envPath, envContent := range environMDs {
			// check if it is already included
			_, ok := agtmd[envPath]
			if ok {
				continue
			}
			// get dir of envPath
			dir := path.Dir(envPath)
			if strings.HasPrefix(fpath, dir) {
				agtmd[envPath] = envContent
			}
		}

		for agtPath, agtContent := range agt.AgentsMD {
			// check if it is already included
			_, ok := agtmd[agtPath]
			if ok {
				continue
			}
			// get dir of agtPath
			dir := path.Dir(agtPath)
			if strings.HasPrefix(fpath, dir) {
				agtmd[agtPath] = agtContent
			}
		}
	}
	// root agent files, debatable if always needed
	for agtPath, agtContent := range agt.AgentsMD {
		if !strings.Contains(agtPath, "/") {
			agtmd[agtPath] = agtContent
		}
	}
	// root project files, always included
	for envPath, envContent := range environMDs {
		if !strings.Contains(envPath, "/") {
			agtmd[envPath] = envContent
		}
	}

	// fmt.Println("USING INSTRUCTION FILES:", slices.Collect(maps.Keys(agtmd)))

	agtmdSorted := make([]config.AgentMD, 0, len(agtmd))
	for p, c := range agtmd {
		agtmdSorted = append(agtmdSorted, config.AgentMD{Path: p, Content: c})
	}
	sort.Slice(agtmdSorted, func(i, j int) bool {
		p1 := agtmdSorted[i].Path
		p2 := agtmdSorted[j].Path

		parts1 := strings.Split(p1, "/")
		parts2 := strings.Split(p2, "/")

		if len(parts1) != len(parts2) {
			return len(parts1) < len(parts2)
		}

		return p1 < p2
	})

	data["agentsMd"] = agtmdSorted

	// All Agent MDs (requested for full list visibility)
	// Merge both agt.AgentsMD and environMDs
	allMDs := make(map[string]string)
	maps.Copy(allMDs, agt.AgentsMD)
	maps.Copy(allMDs, environMDs)

	allAgtmdSorted := make([]config.AgentMD, 0, len(allMDs))
	for p, c := range allMDs {
		allAgtmdSorted = append(allAgtmdSorted, config.AgentMD{Path: p, Content: c})
	}
	sort.Slice(allAgtmdSorted, func(i, j int) bool {
		return allAgtmdSorted[i].Path < allAgtmdSorted[j].Path
	})
	data["allAgentsMd"] = allAgtmdSorted

	stateKeys := slices.Collect(maps.Keys(state))
	slices.Sort(stateKeys)
	data["stateKeys"] = stateKeys

	return data, nil
}

func debugPrintData(data map[string]any) {
	fmt.Println("stateKeys:")
	if keys, ok := data["stateKeys"].([]string); ok {
		for _, k := range keys {
			fmt.Println(" ", k)
		}
	}

	fmt.Println("cacheKeys:")
	if cache, ok := data["cache"].([]KVPair); ok {
		for _, kv := range cache {
			fmt.Println(" ", kv.Key)
		}
	}

	fmt.Println("filesKeys:")
	if files, ok := data["files"].([]KVPair); ok {
		for _, kv := range files {
			fmt.Println(" ", kv.Key)
		}
	}

	fmt.Println("agentKeys:")
	if agtmd, ok := data["agentsMd"].([]config.AgentMD); ok {
		for _, am := range agtmd {
			fmt.Println(" ", am.Path)
		}
	}

	fmt.Println("dataKeys:")
	dataKeys := slices.Collect(maps.Keys(data))
	slices.Sort(dataKeys)
	for _, k := range dataKeys {
		fmt.Println(" ", k)
	}

	fmt.Println("subconscious:", data["subconscious"])
}
