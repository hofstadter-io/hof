package runtime

import "context"

// also move cue agent here

// map of these in Client as well
type Session struct {
	Sid string

	StopFunc context.CancelFunc
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
