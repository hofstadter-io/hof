package ws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

// TODO, we need a good list of message types
// for both the frontend, backend, and where/how they are used

func SetupHandlers(r *runtime.Runtime) {

	// standard fare
	r.Handlers["echo"] = echo
	r.Handlers["hello"] = hello

	// informational handlers
	r.Handlers["requestSync"] = broadcastSync
	r.Handlers["config.reload"] = reloadConfig
	r.Handlers["config.info"] = configInfo
	r.Handlers["models.list"] = modelsList
	r.Handlers["agents.list"] = agentsList

	// chat
	r.Handlers["chat"] = chatUserMessage
	r.Handlers["chat.userMessage"] = chatUserMessage
	r.Handlers["session.cancel"] = sessionCancel

	// sessions
	r.Handlers["session.get"] = sessionGet
	r.Handlers["session.getList"] = sessionList
	r.Handlers["session.create"] = sessionCreate
	r.Handlers["session.delete"] = sessionDelete
	r.Handlers["session.getStateAll"] = sessionGetStateAll
	r.Handlers["session.state.get"] = sessionGetState
	r.Handlers["session.state.put"] = sessionPutState
	r.Handlers["session.state.del"] = sessionDelState

	r.Handlers["session.merge"] = sessionMerge
	r.Handlers["session.tag"] = sessionTag
	r.Handlers["session.push"] = sessionPush
	r.Handlers["session.pull"] = sessionPull
	r.Handlers["session.clone"] = sessionClone
	r.Handlers["session.splice"] = sessionSplice
	// r.Handlers["session.environ.set"] = sessionEnvironSet

	//
	// things we want to track from the frontend
	//
	// TODO, we want to track these on a client basis, so it is available to all agents
	//  then only include some in the data that goes into populating the system prompt
	// r.Handlers["env.info.resp"] = envInfo

}

type EchoPayload struct {
	Text string `json:"text"`
}

type EchoResponsePayload struct {
	ResponseText string `json:"responseText"`
}

func echo(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p EchoPayload
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'echo' payload: %v", err)
		return
	}
	log.Printf("Echoing text: %s", p.Text)

	respPayload := EchoResponsePayload{
		ResponseText: "Server acknowledges: " + p.Text,
	}
	c.Mail("echoResponse", respPayload)
}

type HelloPayload struct {
	Version string `json:"version"`
}

func hello(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p HelloPayload
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'hello' payload: %v", err)
		return
	}
	log.Printf("Hello from client version: %s (Client: %p)", p.Version, c)
}

func broadcastSync(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	// fmt.Println("broadcastSync")
	reloadConfig(r, c, m)
	sessionGet(r, c, m)
	sessionList(r, c, m)
	// sessionFilesysDiff(r, c, m)

	// runtime (runners?)
	// memory
	// artifacts
}

func reloadConfig(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	// todo, this should happen on a per-client/user basis
	var err error
	err = r.ReadConfig()
	if err != nil {
		err = cuetils.ExpandCueError(err)
		c.Mail("config.reload.error", map[string]any{
			"status":        "error",
			"error_message": fmt.Errorf("while reloading config: %w", err),
		})
	}
	configInfo(r, c, m)
}
