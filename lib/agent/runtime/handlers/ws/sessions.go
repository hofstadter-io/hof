package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"slices"
	"time"

	"google.golang.org/adk/session"

	"github.com/google/uuid"
	"github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/agent/runtime/services/environ"
	"github.com/kr/pretty"
)

type SidRequest struct {
	Sid string `json:"sid"`
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
	c.Mail("session.resp.get", S)
	// sessionFilesysDiff(r, c, m)
}

func sessionList(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	// sessions
	sessions, err := r.S.List(r.Ctx, &session.ListRequest{
		AppName: r.AppName,
		UserID:  c.User,
	})
	if err != nil {
		log.Printf("session.getList: %v", err)
		c.Mail("session.list.resp", map[string]string{
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
	var err error
	var payload SessionCreateRequest

	// unpack our payload
	if err := json.Unmarshal(m.Payload, &payload); err != nil {
		log.Printf("Error unmarshaling 'session.create' payload: %v", err)
		return
	}

	fmt.Printf("CREATE SESSION: %#+v\n", pretty.Formatter(payload))

	// initial state
	initialState := make(map[string]any)
	if payload.Title != "" {
		initialState["title"] = payload.Title
	}

	initialState["agent"] = payload.Agent
	initialState["model"] = payload.Model
	initialState["envName"] = payload.EnvName

	pe := payload.Environ
	if pe == nil {
		pe = new(environ.EnvironCreateOptions)
	}

	// maybe attach an environment
	if pe.FromUri == "" && payload.EnvName != "" {
		fmt.Println("searching for env:", payload.EnvName)
		for _, e := range r.Agentic.Environs {
			fmt.Printf(" ? %#+v\n", e)
			if e.Name == payload.EnvName {
				fmt.Println("  MATCH")
				pe.FromUri = "oci://" + e.Spec.From
				break
			}
		}

		env := environ.Client()
		envUri, err := env.Create(pe)
		if err != nil {
			log.Printf("in 'session.create' while creating env: %v", err)
			return
		}
		// will these empty strings get deleted? (vs nil to delete, make sure delete is correct)
		initialState["initEnv"] = pe
		initialState["origEnv"] = string(envUri)
		initialState["currEnv"] = string(envUri)
	}

	// include any client level state
	maps.Copy(initialState, c.State)

	// create our session
	resp, err := r.S.Create(r.Ctx, &session.CreateRequest{
		AppName: r.AppName,
		UserID:  c.User,
		State:   initialState,
	})
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return
	}

	// make sure everyone is notified (just the overall list that most listen to)
	sessionList(r, c, m)
	// sessionFilesysDiff(r, c, m)

	// if focused, tell chat
	if payload.Focus {
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
		c.Mail("session.state.get.resp", map[string]string{
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

	// "create" (put) the session (by using the same Sid)
	err = r.S.AppendEvent(r.Ctx, resp.Session, &session.Event{
		Author:       "user",
		ID:           uuid.NewString(),
		InvocationID: uuid.NewString(),
		Timestamp:    time.Now(),
		Actions: session.EventActions{
			StateDelta: map[string]any{
				s.Key: s.Val,
			},
		},
	})
	if err != nil {
		log.Printf("Error: session.state.put.AppendEvent: %v", err)
		c.Mail("session.state.put.resp", map[string]string{
			"id":    s.Sid,
			"error": err.Error(),
		})
		return
	}

	// fmt.Println("State Set", s.Sid, s.Key, s.Val)
	// fmt.Println("session.state", maps.Collect(resp.Session.State().All()))
}

func sessionDelState(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var s StatePayload
	if err := json.Unmarshal(m.Payload, &s); err != nil {
		log.Printf("Error unmarshaling 'session.state.del.payload': %v", err)
		return
	}
	// lookup session
	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: s.Sid,
	})
	if err != nil {
		log.Printf("Error: session.state.del.getSession: %v", err)
		c.Mail("session.state.del.resp", map[string]string{
			"sid":   s.Sid,
			"error": err.Error(),
		})
		return
	}

	// "create" (put) the session (by using the same Sid)
	err = r.S.AppendEvent(r.Ctx, resp.Session, &session.Event{
		Author:       "user",
		ID:           uuid.NewString(),
		InvocationID: uuid.NewString(),
		Timestamp:    time.Now(),
		Actions: session.EventActions{
			StateDelta: map[string]any{
				s.Key: nil,
			},
		},
	})
	if err != nil {
		log.Printf("Error: session.state.del.AppendEvent: %v", err)
		c.Mail("session.state.del.resp", map[string]string{
			"id":    s.Sid,
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

func sessionFork(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

}

func sessionMerge(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

}

func sessionTag(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

}

func sessionPush(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

}

func sessionPull(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {

}
