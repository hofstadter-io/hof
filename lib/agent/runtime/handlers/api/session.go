package api

import (
	"fmt"
	"iter"
	"maps"
	"net/http"
	"slices"

	"github.com/hofstadter-io/hof/lib/agent/agents"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
	"github.com/labstack/echo/v4"
	"google.golang.org/adk/session"
)

type sessionCloneRequest struct {
	Sid   string `json:"sid"`
	Pos   int    `json:"pos,omitempty"`
	Focus bool   `json:"focus,omitempty"`
}

func (r *Runtime) sessionClone(c echo.Context) error {
	var p sessionCloneRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// 1. Get Session
	sreq := &session.GetRequest{
		AppName:   r.AppName,
		UserID:    "tony",
		SessionID: p.Sid,
	}
	sresp, err := r.S.Get(c.Request().Context(), sreq)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// 2. Clone
	cloned, err := r.S.Clone(c.Request().Context(), sresp.Session)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// splice if pos is non-zero
	if p.Pos > 0 {
		n := cloned.Events().Len()
		if p.Pos < n {
			cloned, err = r.S.Splice(c.Request().Context(), cloned, p.Pos, n-p.Pos, nil)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}
	}

	// build outgoing payload
	S := make(map[string]any)
	S["sid"] = cloned.ID()
	S["state"] = maps.Collect(cloned.State().All())
	S["events"] = slices.Collect(cloned.Events().All())
	S["lastUpdate"] = cloned.LastUpdateTime().UTC()
	S["focus"] = p.Focus

	return c.JSON(http.StatusOK, S)
}

type sessionSpliceRequest struct {
	Sid   string           `json:"sid"`
	Pos   int              `json:"pos"`
	Count int              `json:"count"`
	Fill  []*session.Event `json:"fill"`
}

func (r *Runtime) sessionSplice(c echo.Context) error {
	var p sessionSpliceRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// 1. Get Session
	sreq := &session.GetRequest{
		AppName:   r.AppName,
		UserID:    "tony",
		SessionID: p.Sid,
	}
	sresp, err := r.S.Get(c.Request().Context(), sreq)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// 2. Splice
	spliced, err := r.S.Splice(c.Request().Context(), sresp.Session, p.Pos, p.Count, spliceEvents(p.Fill))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// build outgoing payload
	S := make(map[string]any)
	S["sid"] = spliced.ID()
	S["state"] = maps.Collect(spliced.State().All())
	S["events"] = slices.Collect(spliced.Events().All())
	S["lastUpdate"] = spliced.LastUpdateTime().UTC()

	return c.JSON(http.StatusOK, S)
}

type spliceEvents []*session.Event

func (e spliceEvents) All() iter.Seq[*session.Event] {
	return func(yield func(*session.Event) bool) {
		for _, event := range e {
			if !yield(event) {
				return
			}
		}
	}
}

func (e spliceEvents) Len() int {
	return len(e)
}

func (e spliceEvents) At(i int) *session.Event {
	if i >= 0 && i < len(e) {
		return e[i]
	}
	return nil
}

type promptRenderRequest struct {
	Sid   string `json:"sid"`
	Pos   int    `json:"pos"`
	Agent string `json:"agent"`
}

func (r *Runtime) promptRender(c echo.Context) error {
	var p promptRenderRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// 1. Get Session
	sreq := &session.GetRequest{
		AppName:   r.AppName,
		UserID:    "tony",
		SessionID: p.Sid,
	}
	ctx := c.Request().Context()
	s := r.S
	sresp, err := s.Get(ctx, sreq)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	sess := sresp.Session

	// do we have an env? if yes, get all the agent files for use during instruction generation
	envUri, _ := sess.State().Get("currEnv")
	var environMDs map[string]string
	if envUri != nil {
		environMDs, err = environ.Client().FindAgentFiles(envUri.(string))
		if err != nil {
			fmt.Printf("promptRender.GetAgentFiles.error: %v\n", err)
		}
	}

	// 2. Get Agent Config
	agentName := p.Agent
	if agentName == "" {
		return c.String(http.StatusBadRequest, "agent must be set in request")
		// Try to get agent from state or use a default if available
		// For now, if empty, we might need it passed or found in state
		// v, _ := sess.State().Get("agent")
		// if v != nil {
		// 	agentName = v.(string)
		// }
	}

	agt, err := agents.LoadAgent(r.Agentic, agentName)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// 3. Prepare State
	// TODO: We may need to walk the events backwards and process state changes inversely,
	// from the current state for the session (I don't think it's recorded, only the latest and delta,
	// we can leave this for later, leave a comment where it should go for now
	st := maps.Collect(sess.State().All())

	// 4. Render
	// fmt.Printf("promptRender.render.start: %s\n", agentName)
	prompt, err := agents.RenderInstructionsWithNameAndState(r.Agentic, agt, agentName, st, environMDs)
	if err != nil {
		// fmt.Printf("promptRender.render.error: %v\n", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// TODO calculate tokens here
	// fmt.Printf("promptRender.render.success: %d bytes\n", len(prompt))

	return c.JSON(http.StatusOK, map[string]string{"prompt": prompt})
}
