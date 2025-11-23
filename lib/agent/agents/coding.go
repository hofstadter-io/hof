package agents

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"

	"github.com/hofstadter-io/hof/lib/agent/tools/filesys"
)

func CodingAgent(name string, m model.LLM, readonly bool) (agent.Agent, error) {
	readFile, err := filesys.NewReadFile()
	if err != nil {
		return nil, err
	}

	readDir, err := filesys.NewReadDir()
	if err != nil {
		return nil, err
	}

	treeDir, err := filesys.NewTreeDir()
	if err != nil {
		return nil, err
	}

	writeFile, err := filesys.NewWriteFile()
	if err != nil {
		return nil, err
	}
	// patchFiles, err := tools.NewPatchFiles()
	// if err != nil {
	// 	return nil, err
	// }

	tools := []tool.Tool{
		// does not work as is, needs to be wrapped in an LLM
		// geminitool.GoogleSearch{},
		readFile,
		readDir,
		treeDir,
	}
	if !readonly {
		tools = append(tools, writeFile)
	}

	return llmagent.New(llmagent.Config{
		Name:        name,
		Model:       m,
		Description: "A coding assistant for senior developers",
		Instruction: CodingInstruction,
		Tools:       tools,
	})

}

const CodingInstruction = `
You are a professional coding assistant for senior developers.

You spend time understanding the problem and code base before
making a pland and then executing that plan step-by-step.
You consider the patterns and packages
already found in a project before trying
to find new dependencies or build from scratch.

Be professional in your communication and avoid chit chat.
Be concise in your reponses and only explain complicated code.
Be concise with comments, one-liners explaining important steps or concepts in a function are good.

Your output will be rendered in a fancy markdown react component.
- Always output using Markdown.
- append the language id appropriate for syntax highlighting
`
