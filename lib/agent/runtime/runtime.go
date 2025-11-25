package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/model"
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
	Agentic agents.Config

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
	err = R.ReadConfig()
	if err != nil {
		return fmt.Errorf("while init'n runtime: %w", err)
	}

	err = R.initModels()
	if err != nil {
		return fmt.Errorf("while init'n runtime: %w", err)
	}

	err = R.initServices()
	if err != nil {
		return fmt.Errorf("while init'n runtime: %w", err)
	}

	return nil
}

func (R *Runtime) ReadConfig() error {
	// TODO, load agents from multiple locations
	// 1. user
	// 2. project
	// base on workspaceDir, eventually sent by vs code, or git clone in ephemeral dagger

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
	R.Agentic, err = agents.AgenticCUE(adir, R.Models)
	if err != nil {
		err = cuetils.ExpandCueError(err)
		return fmt.Errorf("while loading AgenticCUE:\n%s", err)
	}

	return nil
}

func (R *Runtime) initModels() (err error) {
	for _, m := range R.Agentic.Models {
		R.Models[m.Name], err = models.Gemini(R.Ctx, m.Id)
		if err != nil {
			return fmt.Errorf("while init'n model %q: %w", m, err)
		}
	}

	return nil
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
