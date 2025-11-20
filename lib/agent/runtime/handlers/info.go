package handlers

import (
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"

	"github.com/hofstadter-io/hof/lib/agent/runtime"
)

type ModelPayload struct {
	Name string
}

func modelToPayload(m model.LLM) ModelPayload {
	return ModelPayload{
		Name: m.Name(),
	}
}

func modelsList(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	fmt.Println("models.list", r.Models)
	models := make(map[string]ModelPayload)
	for k, v := range r.Models {
		models[k] = modelToPayload(v)
	}

	c.Mail("models.list.resp", models)
}

type AgentPayload struct {
	Name string
}

func agentToPayload(a agent.Agent) AgentPayload {
	return AgentPayload{
		Name: a.Name(),
	}
}

func agentsList(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	fmt.Println("agents.list", r.Agents)
	agents := make(map[string]AgentPayload)
	for k, v := range r.Agents {
		agents[k] = agentToPayload(v)
	}

	c.Mail("agents.list.resp", agents)
}
