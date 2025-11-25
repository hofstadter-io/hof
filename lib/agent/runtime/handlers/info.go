package handlers

import (
	"encoding/json"
	"log"

	"github.com/hofstadter-io/hof/lib/agent/runtime"
)

func modelsList(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	c.Mail("models.list.resp", r.Agentic.Models)
}

func agentsList(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	c.Mail("agents.list.resp", r.Agentic.Agents)
}

func envInfo(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	env := make(map[string]any)
	if err := json.Unmarshal(m.Payload, &env); err != nil {
		log.Printf("Error unmarshaling 'envInfo' handler payload: %v", err)
		return
	}
	// fmt.Println("envInfo.input", env)

	// we should probably just structure the object and save it on the client
	// then add some fields when we set it on the session when we chat message

	c.State["env"] = env

	// // ignore if not found
	// esid, ok := env["sid"]
	// if !ok || esid == nil {
	// 	return
	// }

	// sid := esid.(string)
	// // lookup session
	// resp, err := r.S.Get(r.Ctx, &session.GetRequest{
	// 	AppName:   r.AppName,
	// 	UserID:    c.User,
	// 	SessionID: sid,
	// })
	// if err != nil {
	// 	log.Printf("Error: envInfo.get: %v", err)
	// 	c.Mail("session.state.put.resp", map[string]string{
	// 		"sid":   sid,
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// // start building up new state
	// state := maps.Collect(resp.Session.State().All())
	// state["sid"] = sid

	// wsDir := env["workspaceDir"].(string)
	// state["workspaceDir"] = wsDir

	// fmt.Println("envInfo.save", state)

	// // "create" (put) the session (by using the same Sid)
	// _, err = r.S.Create(r.Ctx, &session.CreateRequest{
	// 	AppName:   r.AppName,
	// 	UserID:    c.User,
	// 	SessionID: sid,
	// 	State:     state,
	// })
	// if err != nil {
	// 	log.Printf("Error: envInfo.save: %v", err)
	// 	c.Mail("session.state.put.resp", map[string]string{
	// 		"sid":   sid,
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }
}
