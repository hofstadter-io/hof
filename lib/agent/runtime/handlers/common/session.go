package common

import (
	"google.golang.org/adk/session"

	aruntime "github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/consts"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func SessionGet(r *runtime.Runtime, ar *aruntime.Runtime, sid string) (session.Session, error) {
	resp, err := ar.S.Get(r.Ctx, &session.GetRequest{
		AppName:   ar.AppName,
		UserID:    consts.VEG_DEFAULT_USER,
		SessionID: sid,
	})
	if err != nil {
		return nil, err
	}

	return resp.Session, nil
}

func SessionList(r *runtime.Runtime, ar *aruntime.Runtime) ([]session.Session, error) {
	sessions, err := ar.S.List(r.Ctx, &session.ListRequest{
		AppName: ar.AppName,
		UserID:  consts.VEG_DEFAULT_USER,
	})
	if err != nil {
		return nil, err
	}

	return sessions.Sessions, nil
}
