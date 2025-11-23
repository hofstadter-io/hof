package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"slices"

	"github.com/hofstadter-io/hof/lib/agent/runtime"
	"google.golang.org/adk/session"
)

type SidRequest struct {
	Sid string `json:"sid"`
}

type SessionCreateRequest struct {
	Title string `json:"title,omitempty"`
	Focus bool   `json:"focus,omitempty"`
}

type SessionCreateResponse struct {
	Sid   string `json:"sid"`
	Title string `json:"title,omitempty"`
	Focus bool   `json:"focus,omitempty"`
}

func sessionGet(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

	fmt.Println("sessionGet", string(m.Payload))
	// parse incoming payload
	var p SidRequest
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'session.get' payload: %v", err)
		return
	}

	// lookup session
	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: p.Sid,
	})
	if err != nil {
		log.Printf("session.get: %v", err)
		c.Mail("session.get.resp", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
		return
	}

	// build outgoing payload
	s := resp.Session
	S := make(map[string]any)
	S["sid"] = s.ID()
	S["state"] = maps.Collect(s.State().All())
	S["events"] = slices.Collect(s.Events().All())
	S["lastUpdate"] = s.LastUpdateTime()

	// fmt.Println("mailing sessions", payload)
	c.Mail("session.info", S)
}

func sessionList(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	// sessions
	sessions, err := r.S.List(r.Ctx, &session.ListRequest{
		AppName: r.AppName,
		UserID:  c.User,
	})
	if err != nil {
		log.Printf("session.getList: %v", err)
		c.Mail("session.list", map[string]string{
			"error": err.Error(),
		})
		return
	}

	payload := make([]map[string]any, 0, len(sessions.Sessions))
	for _, s := range sessions.Sessions {
		S := make(map[string]any)
		S["sid"] = s.ID()
		S["state"] = maps.Collect(s.State().All())
		S["events"] = slices.Collect(s.Events().All())
		S["lastUpdate"] = s.LastUpdateTime()
		payload = append(payload, S)
	}
	c.Mail("session.list", payload)
}

func sessionCreate(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p SessionCreateRequest
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'session.delete' payload: %v", err)
		return
	}
	s := make(map[string]any)
	if p.Title != "" {
		s["title"] = p.Title
	}
	maps.Copy(s, c.State)
	resp, err := r.S.Create(r.Ctx, &session.CreateRequest{
		AppName: r.AppName,
		UserID:  c.User,
		State:   s,
	})
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return
	}

	// make sure everyone is notified (just the overall list that most listen to)
	sessionList(r, c, m)

	// if focused, tell chat
	if p.Focus {
		c.Mail("chat.loadSession", map[string]any{
			"sid": resp.Session.ID(),
		})
	}
}

func sessionDelete(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p SidRequest
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'session.delete' payload: %v", err)
		return
	}
	err := r.S.Delete(r.Ctx, &session.DeleteRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: p.Sid,
	})
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return
	}
	// make sure everyone is notified
	// c.broadcastSessions()
	// hacky, but should work the same
	sessionList(r, c, m)
}

type StatePayload struct {
	Sid string `json:"sid"`
	Key string `json:"key"`
	Val any    `json:"val"`
}

func sessionGetStateAll(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

	var s StatePayload
	if err := json.Unmarshal(m.Payload, &s); err != nil {
		log.Printf("Error unmarshaling 'session.getState' payload: %v", err)
		return
	}
	// lookup session
	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: s.Sid,
	})
	if err != nil {
		log.Printf("session.getStateAll: %v", err)
		c.Mail("session.get.resp", map[string]string{
			"id":    s.Sid,
			"error": err.Error(),
		})
		return
	}

	s.Val = maps.Collect(resp.Session.State().All())

	// fmt.Println("mailing sessions", payload)
	c.Mail("session.getStateAll.resp", s)
}

func sessionGetState(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

	var s StatePayload
	if err := json.Unmarshal(m.Payload, &s); err != nil {
		log.Printf("Error unmarshaling 'session.state.get' payload: %v", err)
		return
	}
	c.Mail("session.state.get.req", s)

	// lookup session
	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: s.Sid,
	})
	if err != nil {
		log.Printf("Error: session.state.get.getSession: %v", err)
		c.Mail("session.state.get", map[string]string{
			"id":    s.Sid,
			"error": err.Error(),
		})
		return
	}

	v, err := resp.Session.State().Get(s.Key)
	if err != nil {
		log.Printf("Error: session.state.getState: %v", err)
		c.Mail("session.state.get.resp", map[string]string{
			"id":    s.Sid,
			"error": err.Error(),
		})
	}
	s.Val = v

	// fmt.Println("mailing sessions", payload)
	c.Mail("session.state.get.resp", s)
}

func sessionPutState(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var s StatePayload
	if err := json.Unmarshal(m.Payload, &s); err != nil {
		log.Printf("Error unmarshaling 'session.getState' payload: %v", err)
		return
	}

	fmt.Println("sessionPutState", s)
	// lookup session
	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: s.Sid,
	})
	if err != nil {
		log.Printf("Error: session.state.put.getSession: %v", err)
		c.Mail("session.state.put.resp", map[string]string{
			"id":    s.Sid,
			"error": err.Error(),
		})
		return
	}

	err = resp.Session.State().Set(s.Key, s.Val)
	if err != nil {
		log.Printf("Error: session.state.put.setState: %v", err)
		c.Mail("session.state.put.resp", map[string]string{
			"id":    s.Sid,
			"error": err.Error(),
		})
	}

	// "create" (put) the session (by using the same Sid)
	_, err = r.S.Create(r.Ctx, &session.CreateRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: s.Sid,
		State:     maps.Collect(resp.Session.State().All()),
	})

	// fmt.Println("State Set", s.Sid, s.Key, s.Val)
	// fmt.Println("session.state", maps.Collect(resp.Session.State().All()))
}

func sessionDelState(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var s StatePayload
	if err := json.Unmarshal(m.Payload, &s); err != nil {
		log.Printf("Error unmarshaling 'session.delState' payload: %v", err)
		return
	}
	// lookup session
	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: s.Sid,
	})
	if err != nil {
		log.Printf("Error: session.delState.getSession: %v", err)
		c.Mail("session.delState.resp", map[string]string{
			"sid":   s.Sid,
			"error": err.Error(),
		})
		return
	}

	err = resp.Session.State().Set(s.Key, nil)
	if err != nil {
		log.Printf("Error: session.delState.setState: %v", err)
		c.Mail("session.delState.resp", map[string]string{
			"sid":   s.Sid,
			"error": err.Error(),
		})
	}
}
