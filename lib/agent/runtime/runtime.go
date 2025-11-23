package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/model"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/session/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/hofstadter-io/hof/lib/agent/agents"
	"github.com/hofstadter-io/hof/lib/agent/models"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

// Sqlite driver based on CGO

type Runtime struct {
	AppName string

	Ctx context.Context
	mu  sync.Mutex // To protect clients map
	db  *gorm.DB
	e   *echo.Echo

	// services
	A artifact.Service
	M memory.Service
	S session.Service

	// agentic stuff
	Models  map[string]model.LLM
	Agents  map[string]agent.Agent
	Runners map[string]*runner.Runner

	// clients & comms
	Handlers   map[string]Handler
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

type Handler func(*Runtime, *Client, *Message)

func NewRuntime() (*Runtime, error) {
	R := &Runtime{
		AppName:    "veg",
		Ctx:        context.Background(),
		Models:     make(map[string]model.LLM),
		Agents:     make(map[string]agent.Agent),
		Runners:    make(map[string]*runner.Runner),
		Handlers:   make(map[string]Handler),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	// init components
	err := R.init()
	if err != nil {
		return R, err
	}

	return R, nil
}

func (r *Runtime) ArtifactService() artifact.Service {
	return r.A
}

func (r *Runtime) MemoryService() memory.Service {
	return r.M
}

func (r *Runtime) SessionService() session.Service {
	return r.S
}

func (r *Runtime) handleMessage(c *Client, m *Message) {
	h, ok := r.Handlers[m.Type]
	if ok {
		h(r, c, m)
	}
}

func (R *Runtime) init() (err error) {
	err = R.initModels()
	if err != nil {
		return fmt.Errorf("while init'n runtime: %w", err)
	}
	err = R.initAgents()
	if err != nil {
		return fmt.Errorf("while init'n runtime: %w", err)
	}
	err = R.initServices()
	if err != nil {
		return fmt.Errorf("while init'n runtime: %w", err)
	}
	err = R.initRunners()
	if err != nil {
		return fmt.Errorf("while init'n runtime: %w", err)
	}

	return nil
}

func (R *Runtime) initModels() (err error) {
	ms := []string{
		"gemini-2.5-flash-lite",
		"gemini-2.5-flash",
		"gemini-2.5-pro",
		"gemini-3-pro-preview",
	}

	for _, m := range ms {
		R.Models[m], err = models.Gemini(R.Ctx, m)
		if err != nil {
			return fmt.Errorf("while init'n model %q: %w", m, err)
		}
	}

	return nil
}

func (R *Runtime) initAgents() error {
	// TODO, load agents from multiple locations
	// 1. user
	// 2. project

	// user
	// udir := configdir.LocalConfig("veg", "agents")
	gdir, err := yagu.FindGitRepoAbsPath(".")
	if err != nil {
		return fmt.Errorf("while searching for git root: %w", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("while getting cwd: %w", err)
	}

	bdir := gdir
	if gdir == "" {
		bdir = cwd
	}

	rdir, err := filepath.Rel(cwd, bdir)
	if err != nil {
		return fmt.Errorf("while relativing dir: %w", err)
	}

	adir := filepath.Join(rdir, "./.veg/agents")
	fmt.Println("dirs", gdir, cwd, bdir, rdir, adir)
	// formatting so CUE accepts it (cannot be absolute, cannot be without leading ./ or ../)
	if strings.HasPrefix(adir, ".veg/") {
		adir = "./" + adir
	}

	// Maybe we wait for the above until we hook agents into hof runtime and schemas

	// project, based on cwd, but should probably look for a git root
	as, err := agents.AgenticCUE(adir, R.Models)
	if err != nil {
		err = cuetils.ExpandCueError(err)
		return fmt.Errorf("while loading AgenticCUE:\n%s", err)
	}
	for _, a := range as {
		R.Agents[a.Name()] = a
	}

	// TODO, also load up instruction files
	// AGENTS.md, CLAUDE.md, .github/...
	// and all of their associated skills, subagents, and the like

	return nil

	// // table driven config, will come from CUE eventually (too, some builtins here)
	// type pair struct {
	// 	n string
	// 	m string
	// 	f func(name string, m model.LLM) (agent.Agent, error)
	// }

	// A := []pair{
	// 	{n: "coding-lite", m: "gemini-2.5-flash-lite", f: func(name string, m model.LLM) (agent.Agent, error) {
	// 		return agents.CodingAgent(name, m, false)
	// 	}},
	// 	{n: "coding-fast", m: "gemini-2.5-flash", f: func(name string, m model.LLM) (agent.Agent, error) {
	// 		return agents.CodingAgent(name, m, false)
	// 	}},
	// 	{n: "coding-norm", m: "gemini-2.5-pro", f: func(name string, m model.LLM) (agent.Agent, error) {
	// 		return agents.CodingAgent(name, m, false)
	// 	}},
	// 	{n: "coding-hard", m: "gemini-3-pro-preview", f: func(name string, m model.LLM) (agent.Agent, error) {
	// 		return agents.CodingAgent(name, m, false)
	// 	}},

	// 	{n: "coding-ro-lite", m: "gemini-2.5-flash-lite", f: func(name string, m model.LLM) (agent.Agent, error) {
	// 		return agents.CodingAgent(name, m, true)
	// 	}},
	// 	{n: "coding-ro-fast", m: "gemini-2.5-flash", f: func(name string, m model.LLM) (agent.Agent, error) {
	// 		return agents.CodingAgent(name, m, true)
	// 	}},
	// 	{n: "coding-ro-norm", m: "gemini-2.5-pro", f: func(name string, m model.LLM) (agent.Agent, error) {
	// 		return agents.CodingAgent(name, m, true)
	// 	}},
	// 	{n: "coding-ro-hard", m: "gemini-3-pro-preview", f: func(name string, m model.LLM) (agent.Agent, error) {
	// 		return agents.CodingAgent(name, m, true)
	// 	}},

	// 	{n: "general-lite", m: "gemini-2.5-flash-lite", f: agents.GeneralAgent},
	// 	{n: "general-fast", m: "gemini-2.5-flash", f: agents.GeneralAgent},
	// 	{n: "general-norm", m: "gemini-2.5-pro", f: agents.GeneralAgent},
	// 	{n: "general-hard", m: "gemini-3-pro-preview", f: agents.GeneralAgent},

	// 	{n: "basic-lite", m: "gemini-2.5-flash-lite", f: agents.BasicAgent},
	// 	{n: "basic-fast", m: "gemini-2.5-flash", f: agents.BasicAgent},
	// 	{n: "basic-norm", m: "gemini-2.5-pro", f: agents.BasicAgent},
	// 	{n: "basic-hard", m: "gemini-3-pro-preview", f: agents.BasicAgent},

	// 	{n: "filesys-lite", m: "gemini-2.5-flash-lite", f: agents.ReadWriteFilesysAgent},
	// 	{n: "filesys-fast", m: "gemini-2.5-flash", f: agents.ReadWriteFilesysAgent},
	// 	{n: "filesys-norm", m: "gemini-2.5-pro", f: agents.ReadWriteFilesysAgent},
	// 	{n: "filesys-hard", m: "gemini-3-pro-preview", f: agents.ReadWriteFilesysAgent},
	// }

	// // now create
	// for _, a := range A {
	// 	g, err := a.f(a.m, R.Models[a.m])
	// 	if err != nil {
	// 		return fmt.Errorf("while init'n agent %q: %w", a.n, err)
	// 	}
	// 	R.Agents[a.n] = g
	// }

	// fmt.Println("initAgents.AgenticCUE")

}

func (R *Runtime) initServices() error {
	// open comms to the db
	db, err := gorm.Open(sqlite.Open(".veg/veg.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error creating database session service: %w", err)
	}
	R.db = db

	s, err := database.NewSessionService(db)
	if err != nil {
		return err
	}
	database.AutoMigrate(s)

	R.A = artifact.InMemoryService()
	R.M = memory.InMemoryService()

	R.S = s
	// R.S = session.InMemoryService()

	return nil
}

func (R *Runtime) initRunners() error {
	for k, v := range R.Agents {
		r, err := runner.New(runner.Config{
			AppName:         "veg",
			Agent:           v,
			SessionService:  R.S,
			ArtifactService: R.A,
			MemoryService:   R.M,
		})
		if err != nil {
			return fmt.Errorf("while init'n runner %q: %w", k, err)
		}
		R.Runners[k] = r
	}

	return nil
}
