package agents

import (
	"github.com/hofstadter-io/hof/lib/agent/tools"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
)

func ReadOnlyFilesysAgent(name string, m model.LLM) (agent.Agent, error) {

	readFile, err := tools.NewReadFile()
	if err != nil {
		return nil, err
	}

	readDir, err := tools.NewReadDir()
	if err != nil {
		return nil, err
	}

	// treeDir, err := tools.NewTreeDir()
	// if err != nil {
	// 	return nil, err
	// }

	// globFiles, err := tools.NewGlobFiles()
	// if err != nil {
	// 	return nil, err
	// }

	return llmagent.New(llmagent.Config{
		Name:        name,
		Model:       m,
		Description: "Agent with read-only access to the filesystem",
		Instruction: "You are a helpful assistant that determines what files and directory are relevant and returns their contents and listing as needed.",
		Tools: []tool.Tool{
			readFile,
			readDir,
			// treeDir,
			// globFiles,
		},
	})

}

func ReadWriteFilesysAgent(name string, m model.LLM) (agent.Agent, error) {

	readFile, err := tools.NewReadFile()
	if err != nil {
		return nil, err
	}

	readDir, err := tools.NewReadDir()
	if err != nil {
		return nil, err
	}

	writeFile, err := tools.NewWriteFile()
	if err != nil {
		return nil, err
	}

	return llmagent.New(llmagent.Config{
		Name:        name,
		Model:       m,
		Description: "Agent with read-only access to the filesystem",
		Instruction: "You are a helpful assistant that determines what files and directory are relevant and returns their contents and listing as needed.",
		Tools: []tool.Tool{
			readFile,
			readDir,
			writeFile,
		},
	})

}
