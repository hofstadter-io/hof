package agents

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"

	"github.com/hofstadter-io/hof/lib/agent/tools/filesys"
)

func ReadOnlyFilesysAgent(name string, m model.LLM) (agent.Agent, error) {

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

	grepFiles, err := filesys.NewGrepFiles()
	if err != nil {
		return nil, err
	}

	return llmagent.New(llmagent.Config{
		Name:        name,
		Model:       m,
		Description: "Agent with read-only access to the filesystem",
		Instruction: "You are a helpful assistant that determines what files and directory are relevant and returns their contents and listing as needed.",
		// these involve creating messages
		// BeforeAgentCallbacks: ,
		// AfterAgentCallbacks: ,
		Tools: []tool.Tool{
			readFile,
			readDir,
			treeDir,
			grepFiles,
		},
	})

}

func ReadWriteFilesysAgent(name string, m model.LLM) (agent.Agent, error) {

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

	return llmagent.New(llmagent.Config{
		Name:        name,
		Model:       m,
		Description: "Agent with read-only access to the filesystem",
		Instruction: "You are a helpful assistant that determines what files and directory are relevant and returns their contents and listing as needed.",
		Tools: []tool.Tool{
			readFile,
			readDir,
			treeDir,
			writeFile,
		},
	})

}
