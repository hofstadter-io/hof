package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/model"
	"google.golang.org/adk/session"

	"github.com/hofstadter-io/hof/lib/agent/config"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
)

type Runtime struct {
	AppName string
	S       session.Service
	Agentic *config.Config
}

func (r *Runtime) GetAppName() string {
	return r.AppName
}

func (r *Runtime) GetSessionService() session.Service {
	return r.S
}

func (r *Runtime) GetAgenticConfig() *config.Config {
	return r.Agentic
}

func (r *Runtime) GetArtifactService() artifact.Service {
	return nil
}

func (r *Runtime) GetModels() map[string]model.LLM {
	return nil
}

func (r *Runtime) ReadEnvConfig() error {
	return fmt.Errorf("ReadEnvConfig not implemented for api.Runtime")
}

func (r *Runtime) SetSession(s *common.Session) {
	// api doesn't track active sessions this way
}

func Setup(appName string, e *echo.Echo, s session.Service, a *config.Config) (*Runtime, error) {
	r := &Runtime{
		AppName: appName,
		S:       s,
		Agentic: a,
	}
	//
	// filesystem
	//
	e.POST("/fs/open", fsOpen)
	e.POST("/fs/stat", r.fsStat)
	e.POST("/fs/read", r.fsRead)
	e.POST("/fs/list", r.fsList)
	e.POST("/fs/diff", r.fsDiff)
	e.POST("/fs/write", r.fsWrite)
	e.POST("/fs/delete", r.fsDelete)
	e.POST("/fs/mkdir", r.fsMkdir)
	e.POST("/fs/rename", r.fsRename)
	e.POST("/fs/copy", r.fsCopy)

	e.POST("/env/list", envList)
	e.POST("/prompt/render", r.promptRender)

	e.POST("/session/list", r.sessionList)
	e.POST("/session/get", r.sessionGet)
	e.POST("/session/create", r.sessionCreate)
	e.POST("/session/delete", r.sessionDelete)
	e.POST("/session/clone", r.sessionClone)
	e.POST("/session/splice", r.sessionSplice)
	e.POST("/session/state/get", r.sessionStateGet)
	e.POST("/session/state/put", r.sessionStatePut)
	e.POST("/session/state/del", r.sessionStateDel)

	return r, nil
}
