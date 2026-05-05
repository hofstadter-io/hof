package ws

import (
	"encoding/json"
	"log"
	"maps"
	"slices"

	"google.golang.org/adk/session"
	"github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
)

type SidRequest struct {
	Sid   string `json:"sid"`
	Pos   int    `json:"pos,omitempty"`
	Focus bool   `json:"focus,omitempty"`
}

func sessionGet(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p SidRequest
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'session.get' payload: %v", err)
		return
	}

	s, err := common.SessionGet(r.Ctx, r, c.User, p.Sid)
	if err != nil {
		c.Mail("session.get.resp", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
		return
	}

	S := make(map[string]any)
	S["sid"] = s.ID()
	S["state"] = maps.Collect(s.State().All())
	S["events"] = slices.Collect(s.Events().All())
	S["lastUpdate"] = s.LastUpdateTime().UTC()

	c.Mail("session.info", S)
	c.Mail("session.resp.get", S)
}

func sessionList(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	sessions, err := common.SessionList(r.Ctx, r, c.User)
	if err != nil {
		log.Printf("session.getList: %v", err)
		c.Mail("session.list.resp", map[string]string{
			"error": err.Error(),
		})
		return
	}

	payload := make([]map[string]any, 0, len(sessions))
	for _, s := range sessions {
		S := make(map[string]any)
		S["sid"] = s.ID()
		S["state"] = maps.Collect(s.State().All())
		S["events"] = slices.Collect(s.Events().All())
		S["lastUpdate"] = s.LastUpdateTime().UTC()
		payload = append(payload, S)
	}
	c.Mail("session.list", payload)
	c.Mail("session.list.resp", payload)
}

type SessionCreateRequest struct {
	Title   string                        `json:"title,omitempty"`
	Focus   bool                          `json:"focus,omitempty"`
	Agent   string                        `json:"agent,omitempty"`
	Model   string                        `json:"model,omitempty"`
	EnvName string                        `json:"envName,omitempty"`
	Environ *environ.EnvironCreateOptions `json:"environ,omitempty"`
}

type SessionCreateResponse struct {
	Uri    string `json:"uri"`
	Title  string `json:"title,omitempty"`
	Focus  bool   `json:"focus,omitempty"`
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func sessionCreate(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var payload common.CreatePayload
	if err := json.Unmarshal(m.Payload, &payload); err != nil {
		log.Printf("Error unmarshaling 'session.create' payload: %v", err)
		return
	}
	payload.User = c.User

	sess, err := common.SessionCreate(r.Ctx, r, payload)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		return
	}

	sessionList(r, c, m)

	var focus bool
	var tmp map[string]any
	json.Unmarshal(m.Payload, &tmp)
	if f, ok := tmp["focus"]; ok {
		focus = f.(bool)
	}

	if focus {
		c.Mail("chat.loadSession", map[string]any{
			"sid": sess.ID(),
		})
	}
}

func sessionDelete(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p SidRequest
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'session.delete' payload: %v", err)
		return
	}
	err := common.SessionDel(r.Ctx, r, c.User, p.Sid)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return
	}
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

	sess, err := common.SessionGet(r.Ctx, r, c.User, s.Sid)
	if err != nil {
		log.Printf("session.getStateAll: %v", err)
		c.Mail("session.get.resp", map[string]string{
			"id":    s.Sid,
			"error": err.Error(),
		})
		return
	}

	s.Val = maps.Collect(sess.State().All())
	c.Mail("session.getStateAll.resp", s)
}

func sessionGetState(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var s StatePayload
	if err := json.Unmarshal(m.Payload, &s); err != nil {
		log.Printf("Error unmarshaling 'session.state.get' payload: %v", err)
		return
	}
	c.Mail("session.state.get.req", s)

	val, err := common.SessionStateGet(r.Ctx, r, c.User, s.Sid, s.Key)
	if err != nil {
		log.Printf("Error: session.state.get: %v", err)
		c.Mail("session.state.get.resp", map[string]string{
			"id":    s.Sid,
			"error": err.Error(),
		})
		return
	}
	s.Val = val
	c.Mail("session.state.get.resp", s)
}

func sessionPutState(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var s StatePayload
	if err := json.Unmarshal(m.Payload, &s); err != nil {
		log.Printf("Error unmarshaling 'session.getState' payload: %v", err)
		return
	}

	err := common.SessionStatePut(r.Ctx, r, c.User, s.Sid, s.Key, s.Val)
	if err != nil {
		log.Printf("Error: session.state.put: %v", err)
		c.Mail("session.state.put.resp", map[string]string{
			"id":    s.Sid,
			"error": err.Error(),
		})
		return
	}
}

func sessionDelState(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var s StatePayload
	if err := json.Unmarshal(m.Payload, &s); err != nil {
		log.Printf("Error unmarshaling 'session.state.del.payload': %v", err)
		return
	}
	err := common.SessionStateDel(r.Ctx, r, c.User, s.Sid, s.Key)
	if err != nil {
		log.Printf("Error: session.state.del: %v", err)
		c.Mail("session.state.del.resp", map[string]string{
			"sid":   s.Sid,
			"error": err.Error(),
		})
		return
	}
}

// type SessionFilesysDiffRequest struct {
// 	Sid  string `json:"sid"`
// 	Pos  int    `json:"pos"`
// 	Show bool   `json:"show,omitempty"`
// }

// type SessionFilesysDiffResponse struct {
// 	Sid    string `json:"sid"`
// 	Pos    int    `json:"pos"`
// 	Show   bool   `json:"show,omitempty"`
// 	Status string `json:"status,omitempty"`
// 	Error  string `json:"error,omitempty"`
// }

// func sessionFilesysDiff(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
// 	var p SessionFilesysDiffRequest
// 	if err := json.Unmarshal(m.Payload, &p); err != nil {
// 		log.Printf("Error unmarshaling 'session.diff' payload: %v", err)
// 		return
// 	}
// 	log.Printf("session.diff.payload: %v", p)

// 	// lookup session
// 	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
// 		AppName:   r.AppName,
// 		UserID:    c.User,
// 		SessionID: p.Sid,
// 	})
// 	if err != nil {
// 		log.Printf("session.diff.error: %v", err)
// 		c.Mail("session.diff.resp", map[string]string{
// 			"sid":   p.Sid,
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	// find first and last fs ids
// 	prevUri, _ := resp.Session.State().Get("origEnv")
// 	nextUri, _ := resp.Session.State().Get("currEnv")

// 	if prevUri == nil && nextUri == nil {
// 		return
// 	}

// 	payload, err := environ.Client().DiffDirectory(prevUri.(string), nextUri.(string))
// 	if err != nil {
// 		log.Printf("session.diff.error: %v", err)
// 		c.Mail("session.diff.resp", map[string]string{
// 			"sid":   p.Sid,
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.Mail("session.diff.resp", payload)

// }

func sessionClone(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p SidRequest
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'session.clone' payload: %v", err)
		return
	}

	cloned, err := common.SessionClone(r.Ctx, r, c.User, p.Sid, p.Pos)
	if err != nil {
		log.Printf("session.clone: %v", err)
		c.Mail("session.clone.resp", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
		return
	}

	// build outgoing payload
	S := make(map[string]any)
	S["sid"] = cloned.ID()
	S["state"] = maps.Collect(cloned.State().All())
	S["events"] = slices.Collect(cloned.Events().All())
	S["lastUpdate"] = cloned.LastUpdateTime().UTC()
	S["focus"] = p.Focus

	c.Mail("session.info", S)
	c.Mail("session.clone.resp", S)

	if p.Focus {
		c.Mail("chat.loadSession", map[string]any{
			"sid": cloned.ID(),
		})
	}

	sessionList(r, c, m)
}

func sessionMerge(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

}

func sessionTag(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

}

func sessionPush(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

}

func sessionPull(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

}

type SessionSpliceRequest struct {
	Sid   string           `json:"sid"`
	Pos   int              `json:"pos"`
	Count int              `json:"count"`
	Fill  []*session.Event `json:"fill"`
}

func sessionSplice(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p SessionSpliceRequest
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'session.splice' payload: %v", err)
		return
	}

	spliced, err := common.SessionSplice(r.Ctx, r, c.User, p.Sid, p.Pos, p.Count, p.Fill)
	if err != nil {
		log.Printf("session.splice: %v", err)
		c.Mail("session.splice.resp", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
		return
	}

	S := make(map[string]any)
	S["sid"] = spliced.ID()
	S["state"] = maps.Collect(spliced.State().All())
	S["events"] = slices.Collect(spliced.Events().All())
	S["lastUpdate"] = spliced.LastUpdateTime().UTC()

	c.Mail("session.info", S)
	c.Mail("session.splice.resp", S)

	sessionList(r, c, m)
}
