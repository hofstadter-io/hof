package api

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/adk/session"

	"github.com/hofstadter-io/hof/lib/agent/config"
)

type Runtime struct {
	AppName string
	S       session.Service
	Agentic *config.Config
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
	e.POST("/fs/stat", fsStat)
	e.POST("/fs/read", fsRead)
	e.POST("/fs/list", fsList)
	e.POST("/fs/diff", fsDiff)
	e.POST("/fs/write", fsWrite)
	e.POST("/fs/delete", fsDelete)
	e.POST("/fs/mkdir", fsMkdir)
	e.POST("/fs/rename", fsRename)
	e.POST("/fs/copy", fsCopy)

	e.POST("/env/list", envList)
	e.POST("/prompt/render", r.promptRender)

	e.POST("/session/clone", r.sessionClone)
	e.POST("/session/splice", r.sessionSplice)

	return r, nil
}
