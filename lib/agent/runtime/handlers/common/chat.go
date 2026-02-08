package common

import (
	"context"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"

	"github.com/hofstadter-io/hof/lib/agent/agents"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
	"github.com/hofstadter-io/hof/lib/runtime"
)

type ChatPayload struct {
	User    string `json:"user"`
	Text    string `json:"text"`
	Sid     string `json:"sid"`
	Agent   string `json:"agent"`
	Model   string `json:"model"`
	Environ string `json:"environ"`
}

// TODO, how do we emit events to the TUI and Websocket?
// TODO, how do we tie this to a particular client or user (want more info than username, pass in struct, or lookup based on ID, also authnz eventually)
func SessionChat(r *runtime.Runtime, ar Runtime, p *ChatPayload) (*Session, error) {
	resp, err := ar.GetSessionService().Get(r.Ctx, &session.GetRequest{
		AppName:   ar.GetAppName(),
		UserID:    p.User,
		SessionID: p.Sid,
	})
	if err != nil {
		return nil, fmt.Errorf("common.SessionChat.getSession: %w", err)
	}
	sess := resp.Session

	var environMDs map[string]string

	// do we have an env? if yes, get all the agent files for use during instruction generation
	envUri, err := sess.State().Get("currEnv")
	if err == nil {
		// return nil, fmt.Errorf("common.SessionChat.getCurrEnv: %w", err)
		// do we have agent paths
		if envUri != nil {
			environMDs, err = environ.Client().FindAgentFiles(envUri.(string))
			if err != nil {
				return nil, fmt.Errorf("common.SessionChat.findAgentFiles: %w", err)
			}
			// fmt.Println("FOUND ENVIRON INSTRUCTION FILES:", slices.Collect(maps.Keys(agentMDs)))
		}
	}

	// --- This is how you serialize a typed response ---
	userMsg := genai.NewContentFromText(p.Text, genai.RoleUser)

	// log.Println("userMsg", userMsg)

	// TODO, attach this to the session or client

	// build the agent on demand
	a, err := agents.BuildAgent(ar.GetAgenticConfig(), p.Agent, p.Model, ar.GetModels(), environMDs)
	if err != nil {
		return nil, fmt.Errorf("while building agent %q: %w", p.Agent, err)
	}

	R, err := runner.New(runner.Config{
		AppName:         ar.GetAppName(),
		Agent:           a,
		SessionService:  ar.GetSessionService(),
		ArtifactService: ar.GetArtifactService(),
	})
	if err != nil {
		return nil, fmt.Errorf("while initializing runner for %q: %w", a.Name(), err)
	}

	// setup session chans
	evtChan := make(chan *session.Event, 32)
	errChan := make(chan error, 4)
	chatCtx, chatStop := context.WithCancel(r.Ctx)

	// create agentic session TURN object
	// TODO, can we reuse this stuff? or should it really be per-turn
	// if so, we should probably update the key and add some position info here like len(events)
	s := &Session{
		Sid:       p.Sid,
		EventChan: evtChan,
		ErrorChan: errChan,
		StopFunc:  chatStop,
		Session:   &sess,
	}
	ar.SetSession(s)

	go func() {
		// streamingMode := agent.StreamingModeSSE
		streamingMode := agent.StreamingModeNone
		for event, err := range R.Run(chatCtx, p.User, p.Sid, userMsg, agent.RunConfig{
			StreamingMode: streamingMode,
		}) {
			if err == nil {
				evtChan <- event
			} else {
				errChan <- err
			}
		}
	}()

	return s, nil
}
