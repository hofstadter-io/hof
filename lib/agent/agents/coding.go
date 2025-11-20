package agents

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"

	"github.com/hofstadter-io/hof/lib/agent/tools"
)

func CodingAgent(name string, m model.LLM) (agent.Agent, error) {
	readFile, err := tools.NewReadFile()
	if err != nil {
		return nil, err
	}

	readDir, err := tools.NewReadDir()
	if err != nil {
		return nil, err
	}

	treeDir, err := tools.NewTreeDir()
	if err != nil {
		return nil, err
	}

	writeFile, err := tools.NewWriteFile()
	if err != nil {
		return nil, err
	}
	// patchFiles, err := tools.NewPatchFiles()
	// if err != nil {
	// 	return nil, err
	// }

	return llmagent.New(llmagent.Config{
		Name:        name,
		Model:       m,
		Description: "A coding assistant for senior developers",
		Instruction: CodingInstruction,
		Tools: []tool.Tool{
			// geminitool.GoogleSearch{},
			readFile,
			readDir,
			treeDir,
			writeFile,
			// patchFiles,
		},
	})

}

const CodingInstruction = `
You are a helpful coding assistant for senior developers.
Be concise in your reponses and only explain complicated code.

You consider the patterns and packages
already found in a project before trying
to find new dependencies or build from scratch.

Always output using Markdown
`
