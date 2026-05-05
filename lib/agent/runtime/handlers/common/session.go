package common

import (
	"context"
	"google.golang.org/adk/session"
)

func SessionGet(ctx context.Context, ar Runtime, user, sid string) (session.Session, error) {
	resp, err := ar.GetSessionService().Get(ctx, &session.GetRequest{
		AppName:   ar.GetAppName(),
		UserID:    user,
		SessionID: sid,
	})
	if err != nil {
		return nil, err
	}

	return resp.Session, nil
}

func SessionList(ctx context.Context, ar Runtime, user string) ([]session.Session, error) {
	sessions, err := ar.GetSessionService().List(ctx, &session.ListRequest{
		AppName: ar.GetAppName(),
		UserID:  user,
	})
	if err != nil {
		return nil, err
	}

	return sessions.Sessions, nil
}

func SessionDel(ctx context.Context, ar Runtime, user, sid string) error {
	return ar.GetSessionService().Delete(ctx, &session.DeleteRequest{
		AppName:   ar.GetAppName(),
		UserID:    user,
		SessionID: sid,
	})
}
