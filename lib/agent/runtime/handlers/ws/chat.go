package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	aruntime "github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/runtime"
	"google.golang.org/adk/model"
	"google.golang.org/adk/session"
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
	sess, err := common.SessionGet(r.Ctx, r, c.User, p.Sid)
	if err != nil {
		// log.Printf("session.get: %v", err)
		c.Mail("session.cancel.error", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
		return
	}
	agent, err := sess.State().Get("agent")
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
	err = r.S.AppendEvent(r.Ctx, sess, cEvt)

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
