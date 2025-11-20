package agents

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
)

func BasicAgent(name string, m model.LLM) (agent.Agent, error) {
	return llmagent.New(llmagent.Config{
		Name:        name,
		Model:       m,
		Description: "A general purpose question answering agent",
		Instruction: "You are a helpful assistant that helps the user. You try to be brief in your answers. Format output as Markdown.",
	})
}

func GeneralAgent(name string, m model.LLM) (agent.Agent, error) {

	return llmagent.New(llmagent.Config{
		Name:        name,
		Model:       m,
		Description: "A general purpose question answering agent",
		Instruction: "You are a helpful assistant that helps the user. You try to be brief in your answers. Format output as Markdown. Search and/or ask clarifying questions when you are unsure",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})

}
