package common

import (
	"context"
	"fmt"
	"iter"
	"maps"
	"time"

	"github.com/google/uuid"
	"google.golang.org/adk/session"

	"github.com/hofstadter-io/hof/lib/agent/agents"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
)

func SessionClone(ctx context.Context, ar Runtime, user, sid string, pos int) (session.Session, error) {
	// 1. Get Session
	sess, err := SessionGet(ctx, ar, user, sid)
	if err != nil {
		return nil, err
	}

	// 2. Clone
	cloned, err := ar.GetSessionService().Clone(ctx, sess)
	if err != nil {
		return nil, err
	}

	// splice if pos is non-zero
	if pos > 0 {
		n := cloned.Events().Len()
		if pos < n {
			cloned, err = ar.GetSessionService().Splice(ctx, cloned, pos, n-pos, nil)
			if err != nil {
				return nil, err
			}
		}
	}

	return cloned, nil
}

type SpliceEvents []*session.Event

func (e SpliceEvents) All() iter.Seq[*session.Event] {
	return func(yield func(*session.Event) bool) {
		for _, event := range e {
			if !yield(event) {
				return
			}
		}
	}
}

func (e SpliceEvents) Len() int {
	return len(e)
}

func (e SpliceEvents) At(i int) *session.Event {
	if i >= 0 && i < len(e) {
		return e[i]
	}
	return nil
}

func SessionSplice(ctx context.Context, ar Runtime, user, sid string, pos, count int, fill []*session.Event) (session.Session, error) {
	// 1. Get Session
	sess, err := SessionGet(ctx, ar, user, sid)
	if err != nil {
		return nil, err
	}

	// 2. Splice
	spliced, err := ar.GetSessionService().Splice(ctx, sess, pos, count, SpliceEvents(fill))
	if err != nil {
		return nil, err
	}

	return spliced, nil
}

func SessionStateGet(ctx context.Context, ar Runtime, user, sid, key string) (any, error) {
	sess, err := SessionGet(ctx, ar, user, sid)
	if err != nil {
		return nil, err
	}

	return sess.State().Get(key)
}

func SessionStatePut(ctx context.Context, ar Runtime, user, sid, key string, val any) error {
	sess, err := SessionGet(ctx, ar, user, sid)
	if err != nil {
		return err
	}

	return ar.GetSessionService().AppendEvent(ctx, sess, &session.Event{
		Author:       "user",
		ID:           uuid.NewString(),
		InvocationID: uuid.NewString(),
		Timestamp:    time.Now().UTC(),
		Actions: session.EventActions{
			StateDelta: map[string]any{
				key: val,
			},
		},
	})
}

func SessionStateDel(ctx context.Context, ar Runtime, user, sid, key string) error {
	return SessionStatePut(ctx, ar, user, sid, key, nil)
}

func SessionPromptRender(ctx context.Context, ar Runtime, user, sid, agentName string) (string, error) {
	// 1. Get Session
	sess, err := SessionGet(ctx, ar, user, sid)
	if err != nil {
		return "", err
	}

	// do we have an env?
	envUri, _ := sess.State().Get("currEnv")
	var environMDs map[string]string
	if envUri != nil {
		environMDs, err = environ.Client().FindAgentFiles(envUri.(string))
		if err != nil {
			fmt.Printf("promptRender.GetAgentFiles.error: %v\n", err)
		}
	}

	// 2. Load Agent
	if agentName == "" {
		return "", fmt.Errorf("agent must be set")
	}

	agt, err := agents.LoadAgent(ar.GetAgenticConfig(), agentName)
	if err != nil {
		return "", err
	}

	// 3. Prepare State
	st := maps.Collect(sess.State().All())

	// 4. Render
	prompt, err := agents.RenderInstructionsWithNameAndState(ar.GetAgenticConfig(), agt, agentName, st, environMDs)
	if err != nil {
		return "", err
	}

	return prompt, nil
}
