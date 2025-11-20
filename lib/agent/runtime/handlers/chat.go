package handlers

import (
	"encoding/json"
	"log"

	"github.com/hofstadter-io/hof/lib/agent/runtime"
	"google.golang.org/adk/agent"
	"google.golang.org/genai"
)

type ChatPayload struct {
	Text  string `json:"text"`
	Sid   string `json:"sid"`
	Agent string `json:"agent"`
}

type ChatResponsePayload struct {
	ResponseText string `json:"responseText"`
}

func chatUserMessage(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p ChatPayload
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'chat' payload: %v", err)
		return
	}
	log.Printf("Chatting payload: %#+v", p)

	if p.Agent == "" {
		p.Agent = "general-fast"
	}

	log.Printf("Chatting with %q: %s", p.Agent, p.Text)

	// --- This is how you serialize a typed response ---
	userMsg := genai.NewContentFromText(p.Text, genai.RoleUser)

	log.Println("userMsg", userMsg)

	// streamingMode := agent.StreamingModeSSE
	streamingMode := agent.StreamingModeNone
	for event := range r.Runners[p.Agent].Run(r.Ctx, c.User, p.Sid, userMsg, agent.RunConfig{
		StreamingMode: streamingMode,
	}) {
		log.Printf("chat.event: %v\n", event)
		c.Mail("chat.event", event)
	}
}
