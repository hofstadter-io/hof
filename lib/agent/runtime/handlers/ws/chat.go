package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hofstadter-io/hof/lib/agent/agents"
	aruntime "github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
	"github.com/hofstadter-io/hof/lib/runtime"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

func makeChatUserMessageHandler(r *runtime.Runtime) aruntime.Handler {
	return func(ar *aruntime.Runtime, c *aruntime.Client, m *aruntime.Message) {
		var p common.ChatPayload
		if err := json.Unmarshal(m.Payload, &p); err != nil {
			c.Mail("chat.event.error", map[string]any{
				"agent":         p.Agent,
				"error_message": fmt.Sprintf("Error unmarshaling 'chat' payload: %v", err),
			})
			return
		}

		p.User = c.User

		log.Printf("Chatting payload: %#+v", p)

		s, err := common.SessionChat(r, ar, &p)
		if err != nil {
			log.Printf("chat.msg.error.SessionChat: %v", err)
			c.Mail("chat.event.error", map[string]any{
				"agent":         p.Agent,
				"error_message": fmt.Sprintf("while chatting: %v", err),
			})
			return
		}

		// every time we get an event...
		for e := range s.EventChan {
			// send the message
			c.Mail("chat.event", e)

			// look for any errors
			select {
			case err := <-s.ErrorChan:
				log.Printf("chat.msg.error.SessionChat.loop: %v", err)
				c.Mail("chat.event.error", map[string]any{
					"agent":         p.Agent,
					"error_message": fmt.Sprintf("while chatting: %v", err),
				})
			default:
			}
		}

	}

}

func chatUserMessage(r *aruntime.Runtime, c *aruntime.Client, m *aruntime.Message) {

	var p common.ChatPayload
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		c.Mail("chat.event.error", map[string]any{
			"agent":         p.Agent,
			"error_message": fmt.Sprintf("Error unmarshaling 'chat' payload: %v", err),
		})
		return
	}
	log.Printf("Chatting payload: %#+v", p)

	// lookup session
	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: p.Sid,
	})
	if err != nil {
		log.Printf("chat.msg.error.getSession: %v", err)
		c.Mail("chat.event.error", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
		return
	}
	sess := resp.Session

	var environMDs map[string]string
	// do we have an env? if yes, get all the agent files for use during instruction generation
	envUri, err := sess.State().Get("currEnv")
	if err == nil {
		// log.Printf("chat.msg.error.getCurrEnv: %v", err)
		// c.Mail("chat.event.error", map[string]string{
		// 	"id":    p.Sid,
		// 	"error": err.Error(),
		// })

		// do we have agent paths
		if envUri != nil {
			environMDs, err = environ.Client().FindAgentFiles(envUri.(string))
			if err != nil {
				log.Printf("chat.msg.error.GetAgentFiles: %v", err)
				c.Mail("chat.event.error", map[string]string{
					"id":    p.Sid,
					"error": err.Error(),
				})
			}
			// fmt.Println("FOUND ENVIRON INSTRUCTION FILES:", slices.Collect(maps.Keys(agentMDs)))
		}
	}

	// --- This is how you serialize a typed response ---
	userMsg := genai.NewContentFromText(p.Text, genai.RoleUser)

	log.Println("userMsg", userMsg, c.State)

	// TODO, attach this to the session or client

	// build the agent on demand
	a, err := agents.BuildAgent(r.Agentic, p.Agent, p.Model, r.Models, environMDs)
	if err != nil {
		err = fmt.Errorf("while building agent %q: %w", p.Agent, err)
		fmt.Println("Error:", err)
		c.Mail("chat.event.error", map[string]any{
			"status":        "error",
			"error_message": err.Error(),
		})
		return
	}

	// TODO, also load up instruction files
	// AGENTS.md, CLAUDE.md, .github/...
	// and all of their associated skills, subagents, and the like

	// we construct the runner on demand
	// ...should we also for the agents/tools
	// ...so they can have access to more scope?
	// ...how do we get the write_file to send the contents to vs code instead of writing to disk?
	// ...perhaps through artifacts
	R, err := runner.New(runner.Config{
		AppName:         r.AppName,
		Agent:           a,
		SessionService:  r.S,
		ArtifactService: r.A,
	})
	if err != nil {
		err = fmt.Errorf("while initializing runner for %q: %w", a.Name(), err)
		fmt.Println("Error:", err)
		c.Mail("chat.event.error", map[string]any{
			"status":        "error",
			"error_message": err.Error(),
		})
		return
	}

	// setup subcontext and wait group
	chatCtx, chatStop := context.WithCancel(r.Ctx)

	r.SetSession(&aruntime.Session{
		Sid:      p.Sid,
		StopFunc: chatStop,
	})

	// streamingMode := agent.StreamingModeSSE
	streamingMode := agent.StreamingModeNone
	for event, err := range R.Run(chatCtx, c.User, p.Sid, userMsg, agent.RunConfig{
		StreamingMode: streamingMode,
	}) {
		if err != nil {
			err = fmt.Errorf("while running agent %q: %w", a.Name(), err)
			fmt.Println("ERROR:", err)
			c.Mail("chat.event.error", map[string]any{
				"event":         event,
				"status":        "error",
				"error_message": err.Error(),
			})
			continue
		}
		// log.Printf("chat.event: %v\n", event)
		c.Mail("chat.event", event)
	}

}

func sessionCancel(r *aruntime.Runtime, c *aruntime.Client, m *aruntime.Message) {
	var p SidRequest
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'session.cancel' payload: %v", err)
		return
	}

	s, ok := r.GetSession(p.Sid)
	if !ok && s == nil {
		c.Mail("session.cancel.error", map[string]string{
			"sid":   p.Sid,
			"error": "unknown sid",
		})
		return
	}

	fmt.Println("cancelling:", p.Sid)

	s.StopFunc()
	c.Mail("session.cancel.resp", map[string]string{
		"sid": p.Sid,
	})

	// lookup session
	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: p.Sid,
	})
	if err != nil {
		// log.Printf("session.get: %v", err)
		c.Mail("session.cancel.error", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
		return
	}
	agent, err := resp.Session.State().Get("agent")
	if err != nil {
		// log.Printf("session.get: %v", err)
		c.Mail("session.cancel.error", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
		return
	}

	cEvt := &session.Event{
		LLMResponse: model.LLMResponse{
			FinishReason: "OTHER",
			TurnComplete: true,
			Interrupted:  true,
		},
		Author:       agent.(string),
		ID:           uuid.NewString(),
		InvocationID: uuid.NewString(),
		Timestamp:    time.Now().UTC(),
		Actions: session.EventActions{
			StateDelta: map[string]any{
				"canceled": true,
			},
		},
	}
	fmt.Println("saving:", p.Sid, cEvt)
	// "create" (put) the session (by using the same Sid)
	err = r.S.AppendEvent(r.Ctx, resp.Session, cEvt)

	if err != nil {
		// log.Printf("session.get: %v", err)
		c.Mail("session.cancel.error", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
		return
	}

	c.Mail("session.cancel.resp", map[string]string{
		"sid": p.Sid,
	})
}
