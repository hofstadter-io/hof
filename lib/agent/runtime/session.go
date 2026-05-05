package runtime

import (
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
)

// also move cue agent here

// map of these in Client as well

func (R *Runtime) GetSession(sid string) (*common.Session, bool) {
	R.sessionsMx.RLock()
	defer R.sessionsMx.RUnlock()

	s, ok := R.sessions[sid]
	return s, ok
}

func (R *Runtime) SetSession(s *common.Session) {
	R.sessionsMx.Lock()
	defer R.sessionsMx.Unlock()
	R.sessions[s.Sid] = s
}
