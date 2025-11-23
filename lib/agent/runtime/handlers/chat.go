package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hofstadter-io/hof/lib/agent/runtime"
	"google.golang.org/adk/agent"
	"google.golang.org/genai"
)

type ChatPayload struct {
	Text  string `json:"text"`
	Sid   string `json:"sid"`
	Agent string `json:"agent"`
	Model string `json:"model"`
}

type ChatResponsePayload struct {
	ResponseText string `json:"responseText"`
}

func chatUserMessage(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p ChatPayload
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		c.Mail("chat.event.error", map[string]any{
			"agent":         p.Agent,
			"error_message": fmt.Sprintf("Error unmarshaling 'chat' payload: %v", err),
		})
		return
	}
	log.Printf("Chatting payload: %#+v", p)

	// --- This is how you serialize a typed response ---
	userMsg := genai.NewContentFromText(p.Text, genai.RoleUser)

	log.Println("userMsg", userMsg, c.State)

	// do we construct agents/tools on demand, so they can have access to more scope?
	// how do we get the write_file to send the contents to vs code instead of writing to disk?

	// TODO, be better about validating inputs
	R, ok := r.Runners[p.Agent]
	if !ok {
		c.Mail("chat.event.error", map[string]any{
			"agent":         p.Agent,
			"error_message": "agent not found",
		})
		return
	}

	// streamingMode := agent.StreamingModeSSE
	streamingMode := agent.StreamingModeNone
	for event := range R.Run(r.Ctx, c.User, p.Sid, userMsg, agent.RunConfig{
		StreamingMode: streamingMode,
	}) {
		log.Printf("chat.event: %v\n", event)
		c.Mail("chat.event", event)
	}

}
