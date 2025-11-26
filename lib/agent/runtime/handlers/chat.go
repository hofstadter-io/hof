package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hofstadter-io/hof/lib/agent/agents"
	"github.com/hofstadter-io/hof/lib/agent/runtime"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
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

	// TODO, attach this to the session or client

	// build the agent on demand
	a, err := agents.BuildAgent(r.Agentic, p.Agent, p.Model, r.Models)
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
		AppName:         "veg",
		Agent:           a,
		SessionService:  r.S,
		ArtifactService: r.A,
		MemoryService:   r.M,
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

	// streamingMode := agent.StreamingModeSSE
	streamingMode := agent.StreamingModeNone
	for event, err := range R.Run(r.Ctx, c.User, p.Sid, userMsg, agent.RunConfig{
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
		log.Printf("chat.event: %v\n", event)
		c.Mail("chat.event", event)
	}

}
