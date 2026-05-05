package runtime

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"cuelang.org/go/cue"
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
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
	vegsession "github.com/hofstadter-io/hof/lib/agent/services/session"
	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/consts"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/env"
	hofruntime "github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/templates"
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

	// hof runtime
	HofRuntime *hofruntime.Runtime

	// Copying(read-only) in temporarily(?) until more of the things here get lifted
	// it at least lets us start refactoring code here around the top-level runtime and CUE fabric
	Envs     []*env.Env
	Agentics []*agent.Agentic

	// clients & comms
	// TODO, this is stuff we should move up and support multiple subsystems with
	ApiRuntime *api.Runtime
	Handlers   map[string]Handler
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client

	sessions   map[string]*common.Session
	sessionsMx sync.RWMutex
}

type Handler func(*Runtime, *Client, *Message)

func NewRuntime(
	hr *hofruntime.Runtime,
) (*Runtime, error) {
	ctx := context.Background()

	R := &Runtime{
		AppName:    "veg",
		Ctx:        ctx,
		db:         hr.DB,
		Envs:       hr.Envs,
		Agentics:   hr.Agentics,
		HofRuntime: hr,
		Models:     make(map[string]model.LLM),
		Handlers:   make(map[string]Handler),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		sessions:   make(map[string]*common.Session),
	}

	// init components
	err := R.init()
	if err != nil {
		return R, err
	}

	return R, nil
}

func (r *Runtime) GetAppName() string {
	return r.AppName
}

func (r *Runtime) GetSessionService() session.Service {
	return r.S
}

func (r *Runtime) GetAgenticConfig() *agentconfig.Config {
	return r.Agentic
}

func (r *Runtime) GetArtifactService() artifact.Service {
	return r.A
}

func (r *Runtime) GetModels() map[string]model.LLM {
	return r.Models
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

			switch m.Provider {

			case "vertex":
				R.Models[m.Name], err = models.Vertex(R.Ctx, m.Id)
				if err != nil {
					return fmt.Errorf("while init'n model %q: %w", m, err)
				}

			case "openai":
				R.Models[m.Name], err = models.OpenAI(R.Ctx, m.Id, m.BaseURL)
				if err != nil {
					return fmt.Errorf("while init'n model %q: %w", m, err)
				}

			default:
				return fmt.Errorf("while init'n model %q, unknown provider: %q", m, m.Provider)

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

	apiRuntime, err := api.Setup(r.AppName, e, r.S, r.Agentic)
	if err != nil {
		return err
	}

	// save & return
	r.ApiRuntime = apiRuntime

	r.e = e
	return nil
}

func (r *Runtime) BackfillAgentic() error {
	cfg := agentconfig.NewConfig()

	// fmt.Println("Agentics", len(r.Agentics))
	for _, a := range r.Agentics {
		switch a.Hof.Agentic.Kind {
		case "agent":
			var m agentconfig.Agent
			err := a.Value.Decode(&m)
			if err != nil {
				return err
			}
			cfg.Agents[a.Hof.Agentic.Name] = m
		case "model":
			var m agentconfig.Model
			err := a.Value.Decode(&m)
			if err != nil {
				return err
			}
			cfg.Models[a.Hof.Agentic.Name] = m
		case "tool":
			var m agentconfig.Tool
			err := a.Value.Decode(&m)
			if err != nil {
				return err
			}
			cfg.Tools[a.Hof.Agentic.Name] = m
		case "environ":
			var m agentconfig.Environ
			err := a.Value.Decode(&m)
			if err != nil {
				return err
			}
			cfg.Environs[a.Hof.Agentic.Name] = m

		case "embed":
			switch a.Value.IncompleteKind() {
			// a path->content map
			case cue.StructKind:
				m := make(map[string]string)
				err := a.Value.Decode(&m)
				if err != nil {
					return err
				}
				for k, v := range m {
					t, err := templates.CreateFromString(k, v, templates.Delims{})
					if err != nil {
						fmt.Println("ERROR.RenderInstructions.Create", err)
						return err
					}
					cfg.Templates[k] = t
				}
			}

		}
	}
	prepareTemplates(cfg)

	r.Agentic = cfg
	// TODO, this needs to be per client / session / workspace
	r.ApiRuntime.Agentic = cfg

	// fmt.Printf("%#+v\n", pretty.Formatter(cfg))

	return nil
}

// TODO, we still need to do this, through probably on a per-request basis, or at least with watch/config changes
func prepareTemplates(cfg *agentconfig.Config) error {

	cwd, _ := os.Getwd()
	// todo, also put this on the Session
	dir := filepath.Join(cwd, cfg.EmbedDir)
	glob := filepath.Join(dir, "**/*.*")
	cfg.Templates = templates.NewTemplateMap()
	// fmt.Printf("found %d templates in %q\n", len(config.Templates), dir)
	err := cfg.Templates.ImportFromFolder(glob, dir, templates.Delims{}, nil)
	if err != nil {
		return fmt.Errorf("while loading instruction templates (%s,%s): %w", cwd, cfg.EmbedDir, err)
	}
	// fmt.Printf("found %d templates in %s\n", len(config.Templates), dir)

	for _, T1 := range cfg.Templates {
		for _, T2 := range cfg.Templates {
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
