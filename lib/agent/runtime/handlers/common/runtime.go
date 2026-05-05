package common

import (
	"context"

	"google.golang.org/adk/artifact"
	"google.golang.org/adk/model"
	"google.golang.org/adk/session"

	"github.com/hofstadter-io/hof/lib/agent/config"
)

type Session struct {
	Sid string

	EventChan chan *session.Event
	ErrorChan chan error
	StopFunc  context.CancelFunc

	Agentic config.Config

	Session *session.Session
}

type Runtime interface {
	GetAppName() string
	GetSessionService() session.Service
	GetAgenticConfig() *config.Config
	GetArtifactService() artifact.Service
	GetModels() map[string]model.LLM
	ReadEnvConfig() error
	SetSession(s *Session)
}
