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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/hofstadter-io/hof/lib/agent/agents"
	"github.com/hofstadter-io/hof/lib/agent/models"
	"github.com/hofstadter-io/hof/lib/agent/runtime/services/environ"
	vegsession "github.com/hofstadter-io/hof/lib/agent/runtime/services/session"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

// TODO, make these env vars
const CONFIG_PATH = `.veg`
const DATA_PATH = `.veg/data`

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
	Agentic agents.Config

	// clients & comms
	Handlers   map[string]Handler
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

type Handler func(*Runtime, *Client, *Message)

func NewRuntime() (*Runtime, error) {
	ctx := context.Background()

	R := &Runtime{
		AppName:    "veg",
		Ctx:        ctx,
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
		return fmt.Errorf("while reading config: %w", err)
	}

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

	adir := filepath.Join(rdir, CONFIG_PATH)
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
	dia := sqlite.Open(filepath.Join(DATA_PATH, "veg.db"))
	db, err := gorm.Open(dia, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error creating database session service: %w", err)
	}
	R.db = db

	// environment management
	err = environ.Initialize(R.Ctx, db)
	if err != nil {
		return fmt.Errorf("while initializing Runtime.EnvironService")
	}

	// session management
	s, err := vegsession.NewSessionServiceGorm(db)
	if err != nil {
		return fmt.Errorf("while initializing Runtime.SessionService")
	}
	vegsession.AutoMigrate(s)
	R.S = s

	// artifacts
	R.A, err = artifact.FilesystemService(filepath.Join(DATA_PATH, "artifacts"))
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

	//
	// filesystem
	//
	e.POST("/fs/open", fsOpen)
	e.POST("/fs/stat", fsStat)
	e.POST("/fs/read", r.fsRead)
	e.POST("/fs/list", r.fsList)
	e.POST("/fs/diff", r.fsDiff)
	e.POST("/fs/write", r.fsWrite)
	e.POST("/fs/delete", r.fsDelete)

	e.POST("/env/list", r.envList)

	// save & return
	r.e = e
	return nil
}
