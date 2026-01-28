package runtime

import (
	"context"

	"google.golang.org/adk/session"

	"github.com/hofstadter-io/hof/lib/agent/config"
)

// also move cue agent here

// map of these in Client as well
type Session struct {
	Sid string

	EventChan chan *session.Event
	ErrorChan chan error
	StopFunc  context.CancelFunc

	Agentic config.Config

	Session *session.Session
}

func (R *Runtime) GetSession(sid string) (*Session, bool) {
	R.sessionsMx.RLock()
	defer R.sessionsMx.RUnlock()

	s, ok := R.sessions[sid]
	return s, ok
}

func (R *Runtime) SetSession(s *Session) {
	R.sessionsMx.Lock()
	defer R.sessionsMx.Unlock()
	R.sessions[s.Sid] = s
}
