package handlers

import (
	"encoding/json"
	"log"

	"github.com/hofstadter-io/hof/lib/agent/runtime"
)

type IdRequest struct {
	ID string `json:"id"`
}

// TODO, we need a good list of message types
// for both the frontend, backend, and where/how they are used

func SetupHandlers(r *runtime.Runtime) {

	// standard fare
	r.Handlers["echo"] = echo
	r.Handlers["hello"] = hello

	// informational handlers
	r.Handlers["requestSync"] = broadcastSync
	r.Handlers["models.list"] = modelsList
	r.Handlers["agents.list"] = agentsList

	// chat
	r.Handlers["chat"] = chatUserMessage
	r.Handlers["chat.userMessage"] = chatUserMessage

	// sessions
	r.Handlers["session.get"] = sessionGet
	r.Handlers["session.getList"] = sessionList
	r.Handlers["session.create"] = sessionCreate
	r.Handlers["session.delete"] = sessionDelete
	r.Handlers["session.getStateAll"] = sessionGetStateAll
	r.Handlers["session.state.get"] = sessionGetState
	r.Handlers["session.state.put"] = sessionPutState

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
	sessionList(r, c, m)
	modelsList(r, c, m)
	agentsList(r, c, m)

	// runtime (runners?)
	// memory
	// artifacts
}
