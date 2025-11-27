package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"

	"dagger.io/dagger"
	"google.golang.org/adk/session"

	"github.com/hofstadter-io/hof/lib/agent/runtime"
	vegdagger "github.com/hofstadter-io/hof/lib/agent/runtime/dagger"
)

type SidRequest struct {
	Sid string `json:"sid"`
}

type SessionCreateRequest struct {
	Title   string `json:"title,omitempty"`
	Dir     string `json:"dir,omitempty"`
	Runtime string `json:"runtime,omitempty"`
	Focus   bool   `json:"focus,omitempty"`
}

type SessionCreateResponse struct {
	Sid    string `json:"sid"`
	Title  string `json:"title,omitempty"`
	Focus  bool   `json:"focus,omitempty"`
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
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
	var err error

	var payload SessionCreateRequest
	if err := json.Unmarshal(m.Payload, &payload); err != nil {
		log.Printf("Error unmarshaling 'session.create' payload: %v", err)
		return
	}

	// initial state
	initialState := make(map[string]any)
	if payload.Title != "" {
		initialState["title"] = payload.Title
	}
	dir := payload.Dir
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			log.Printf("Error in 'session.create' while getting cwd: %v", err)
			return
		}
	}
	fmt.Println("Initializing session with dir", dir)
	// TODO, From container if in config
	d := r.Dagger.Host().Directory(dir, dagger.HostDirectoryOpts{
		Gitignore: true,
	})
	id, err := d.ID(r.Ctx)
	if err != nil {
		log.Printf("Error in 'session.create' while loading dir into dagger: %v", err)
		return
	}
	initialState["origfs"] = string(id)
	initialState["dagger"] = string(id)

	maps.Copy(initialState, c.State)
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

type SessionFilesysDiffRequest struct {
	Sid string `json:"sid"`
	Pos int    `json:"pos"`
}

type SessionFilesysDiffResponse struct {
	Sid    string `json:"sid"`
	Pos    int    `json:"pos"`
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func sessionFilesysDiff(r *runtime.Runtime, c *runtime.Client, m *runtime.Message) {
	var p SessionFilesysDiffRequest
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		log.Printf("Error unmarshaling 'session.diff' payload: %v", err)
		return
	}
	log.Printf("session.diff.payload: %v", p)

	// lookup session
	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
		AppName:   r.AppName,
		UserID:    c.User,
		SessionID: p.Sid,
	})
	if err != nil {
		log.Printf("session.diff: %v", err)
		c.Mail("session.diff.resp", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
	}

	// get dagger handle
	dag, _ := vegdagger.Get(r.Ctx)

	// find first and last fs ids
	origId, _ := resp.Session.State().Get("origfs")
	dagId, _ := resp.Session.State().Get("dagger")

	// walk from 0->pos
	if p.Pos > 0 {
		events := slices.Collect(resp.Session.Events().All())
		posId := origId
		for i := 0; i < p.Pos && i < len(events); i++ {
			event := events[i]
			if did, ok := event.Actions.StateDelta["dagger"]; ok && did != "" && did != posId {
				posId = did
			}
		}
		dagId = posId
	}

	// now get our dirs
	origDir := dag.LoadDirectoryFromID(dagger.DirectoryID(origId.(string)))
	dagDir := dag.LoadDirectoryFromID(dagger.DirectoryID(dagId.(string)))
	changes := dagDir.Changes(origDir)
	// fmt.Println("session.diff.debug", origId, dagId, maps.Collect(resp.Session.State().All()))

	addpaths, err := changes.AddedPaths(r.Ctx)
	if err != nil {
		log.Printf("session.diff.resp: %v", err)
		c.Mail("session.diff.resp", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
	}

	modpaths, err := changes.ModifiedPaths(r.Ctx)
	if err != nil {
		log.Printf("session.diff.resp: %v", err)
		c.Mail("session.diff.resp", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
	}

	delpaths, err := changes.RemovedPaths(r.Ctx)
	if err != nil {
		log.Printf("session.diff.resp: %v", err)
		c.Mail("session.diff.resp", map[string]string{
			"sid":   p.Sid,
			"error": err.Error(),
		})
	}

	pfile := changes.AsPatch()
	patch, err := pfile.Contents(r.Ctx)

	c.Mail("session.diff.resp", map[string]any{
		"sid":      p.Sid,
		"addpaths": addpaths,
		"modpaths": modpaths,
		"delpaths": delpaths,
		"patch":    patch,
	})

}
