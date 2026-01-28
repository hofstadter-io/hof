package runtime

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/model"
	"google.golang.org/adk/session"
	"gorm.io/gorm"

	"github.com/hofstadter-io/hof/lib/agent"
	"github.com/hofstadter-io/hof/lib/agent/agents"
	agentconfig "github.com/hofstadter-io/hof/lib/agent/config"
	"github.com/hofstadter-io/hof/lib/agent/models"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/api"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
	vegsession "github.com/hofstadter-io/hof/lib/agent/services/session"
	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/consts"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/yagu"
)

// TODO, make these env vars

type Runtime struct {
	AppName string

	Ctx context.Context
	mu  sync.Mutex // To protect clients map among other things
	db  *gorm.DB
	e   *echo.Echo

	// services
	A artifact.Service
	S session.Service

	// agentic stuff
	Models  map[string]model.LLM
	Agentic *agentconfig.Config

	// Copying(read-only) in temporarily(?) until more of the things here get lifted
	// it at least lets us start refactoring code here around the top-level runtime and CUE fabric
	Envs     []*env.Env
	Agentics []*agent.Agentic

	// clients & comms
	// TODO, this is stuff we should move up and support multiple subsystems with
	Handlers   map[string]Handler
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client

	sessions   map[string]*Session
	sessionsMx sync.RWMutex
}

type Handler func(*Runtime, *Client, *Message)

func NewRuntime(
	db *gorm.DB,
	envs []*env.Env,
	agentics []*agent.Agentic,
) (*Runtime, error) {
	ctx := context.Background()

	R := &Runtime{
		AppName:    "veg",
		Ctx:        ctx,
		db:         db,
		Envs:       envs,
		Agentics:   agentics,
		Models:     make(map[string]model.LLM),
		Handlers:   make(map[string]Handler),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		sessions:   make(map[string]*Session),
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
	// err = R.ReadConfig()
	// if err != nil {
	// 	return fmt.Errorf("while reading config: %w", err)
	// }

	err = R.initModels()
	if err != nil {
		return fmt.Errorf("while init'n models: %w", err)
	}

	err = R.initServices()
	if err != nil {
		return fmt.Errorf("while init'n services: %w", err)
	}

	err = R.initServer()
	if err != nil {
		return fmt.Errorf("while init'n server: %w", err)
	}

	return nil
}

// this needs to be updated to read out of a fs / env
func (R *Runtime) ReadEnvConfig() error {
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

	adir := filepath.Join(rdir, consts.VEG_REPO_LOCAL_PATH)
	// fmt.Println("dirs", gdir, cwd, bdir, rdir, adir)
	// formatting so CUE accepts it (cannot be absolute, cannot be without leading ./ or ../)
	if strings.HasPrefix(adir, ".veg/") {
		adir = "./" + adir
	}

	// Maybe we wait for the above until we hook agents into hof runtime and schemas

	// project, based on cwd, but should probably look for a git root

	R.Agentic, err = agents.OldAgenticCUE(adir, R.Models)
	if err != nil {
		err = cuetils.ExpandCueError(err)
		return fmt.Errorf("while loading AgenticCUE:\n%s", err)
	}

	return nil
}

func (R *Runtime) initModels() (err error) {
	for _, a := range R.Agentics {
		if a.Hof.Agentic.Kind == "model" {
			var m agentconfig.Model
			err := a.Value.Decode(&m)
			if err != nil {
				return fmt.Errorf("while decoding'n model %q: %w", a.Value, err)
			}
			R.Models[m.Name], err = models.Gemini(R.Ctx, m.Id)
			if err != nil {
				return fmt.Errorf("while init'n model %q: %w", m, err)
			}
		}
	}

	return nil
}

func (R *Runtime) initServices() error {

	// VEG|RENAME: make this a multi-tier lookup and unify system
	// generally for all the subsystems

	// environment management
	err := environ.Initialize(R.Ctx, R.db)
	if err != nil {
		return fmt.Errorf("while initializing Runtime.EnvironService")
	}

	// session management
	s, err := vegsession.NewSessionServiceGorm(R.db)
	if err != nil {
		return fmt.Errorf("while initializing Runtime.SessionService")
	}
	vegsession.AutoMigrate(s)
	R.S = s

	// artifacts
	R.A, err = artifact.FilesystemService(filepath.Join(config.Veg.UserDataDir, "artifacts"))
	if err != nil {
		return fmt.Errorf("while initializing Runtime.ArtifactService")
	}

	return nil
}

func (r *Runtime) initServer() error {
	e := echo.New()
	e.HideBanner = true

	// middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// routes
	e.GET("/", r.serveWs)

	e.GET("/alive", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// TODO metrics & otel

	api.Setup(r.AppName, e, r.S)

	// save & return
	r.e = e
	return nil
}
