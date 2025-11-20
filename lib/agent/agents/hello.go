package agents

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"

	"github.com/hofstadter-io/hof/lib/agent/tools"
)

func NewHelloAgent(m model.LLM) (agent.Agent, error) {
	sumTool, err := tools.NewSummer()
	if err != nil {
		return nil, err
	}

	subTool, err := tools.NewSubber()
	if err != nil {
		return nil, err
	}

	return llmagent.New(llmagent.Config{
		Name:        "hello_time_agent",
		Model:       m,
		Description: "Tells the current time in a specified city.",
		Instruction: "You are a helpful assistant that tells the current time in a city.",
		Tools: []tool.Tool{
			// geminitool.GoogleSearch{},
			// loadartifactstool.New(),
			sumTool,
			subTool,
		},
	})

}
