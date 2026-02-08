package ws

import (
	"encoding/json"
	"fmt"
	"log"

	aruntime "github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func SetupHandlers(r *runtime.Runtime, ar *aruntime.Runtime) {

	// standard fare
	ar.Handlers["echo"] = echo
	ar.Handlers["hello"] = hello

	// informational handlers
	ar.Handlers["requestSync"] = broadcastSync
	ar.Handlers["config.reload"] = reloadEnvConfig
	ar.Handlers["config.info"] = configInfo
	ar.Handlers["models.list"] = modelsList
	ar.Handlers["agents.list"] = agentsList

	// chat
	ar.Handlers["chat"] = makeChatUserMessageHandler(r)
	ar.Handlers["chat.userMessage"] = makeChatUserMessageHandler(r)
	ar.Handlers["session.cancel"] = sessionCancel

	// sessions
	ar.Handlers["session.get"] = sessionGet
	ar.Handlers["session.getList"] = sessionList
	ar.Handlers["session.create"] = sessionCreate
	ar.Handlers["session.delete"] = sessionDelete
	ar.Handlers["session.getStateAll"] = sessionGetStateAll
	ar.Handlers["session.state.get"] = sessionGetState
	ar.Handlers["session.state.put"] = sessionPutState
	ar.Handlers["session.state.del"] = sessionDelState

	ar.Handlers["session.merge"] = sessionMerge
	ar.Handlers["session.tag"] = sessionTag
	ar.Handlers["session.push"] = sessionPush
	ar.Handlers["session.pull"] = sessionPull
	ar.Handlers["session.clone"] = sessionClone
	ar.Handlers["session.splice"] = sessionSplice
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

func echo(r *aruntime.Runtime, c *aruntime.Client, m *aruntime.Message) {
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

func hello(ar *aruntime.Runtime, c *aruntime.Client, m *aruntime.Message) {
	var p HelloPayload
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'hello' payload: %v", err)
		return
	}
	log.Printf("Hello from client version: %s (Client: %p)", p.Version, c)
}

func broadcastSync(ar *aruntime.Runtime, c *aruntime.Client, m *aruntime.Message) {
	// fmt.Println("broadcastSync")
	// reloadEnvConfig(r, c, m)
	sessionGet(ar, c, m)
	sessionList(ar, c, m)
	// sessionFilesysDiff(r, c, m)

	// runtime (runners?)
	// memory
	// artifacts
}

func reloadEnvConfig(ar *aruntime.Runtime, c *aruntime.Client, m *aruntime.Message) {
	// todo, this should happen on a per-client/user basis
	var err error
	err = common.ReloadConfig(ar)
	if err != nil {
		err = cuetils.ExpandCueError(err)
		c.Mail("config.reload.error", map[string]any{
			"status":        "error",
			"error_message": fmt.Errorf("while reloading config: %w", err),
		})
	}
	configInfo(ar, c, m)
}
