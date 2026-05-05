package api

import (
	"maps"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"google.golang.org/adk/session"

	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/consts"
)

type sessionCloneRequest struct {
	Sid   string `json:"sid"`
	Pos   int    `json:"pos,omitempty"`
	Focus bool   `json:"focus,omitempty"`
}

type SidRequest struct {
	Sid string `json:"sid"`
}

func (r *Runtime) sessionClone(c echo.Context) error {
	var p sessionCloneRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	cloned, err := common.SessionClone(c.Request().Context(), r, user, p.Sid, p.Pos)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
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

func (r *Runtime) sessionCreate(c echo.Context) error {
	var p common.CreatePayload
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = p.User
	}
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}
	p.User = user

	sess, err := common.SessionCreate(c.Request().Context(), r, p)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	S := make(map[string]any)
	S["sid"] = sess.ID()
	S["state"] = maps.Collect(sess.State().All())
	S["events"] = slices.Collect(sess.Events().All())
	S["lastUpdate"] = sess.LastUpdateTime().UTC()

	return c.JSON(http.StatusOK, S)
}

func (r *Runtime) sessionGet(c echo.Context) error {
	sid := c.QueryParam("sid")
	if sid == "" {
		var p SidRequest
		err := c.Bind(&p)
		if err == nil && p.Sid != "" {
			sid = p.Sid
		}
	}
	if sid == "" {
		return c.String(http.StatusBadRequest, "missing sid")
	}

	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	sess, err := common.SessionGet(c.Request().Context(), r, user, sid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	S := make(map[string]any)
	S["sid"] = sess.ID()
	S["state"] = maps.Collect(sess.State().All())
	S["events"] = slices.Collect(sess.Events().All())
	S["lastUpdate"] = sess.LastUpdateTime().UTC()

	return c.JSON(http.StatusOK, S)
}

func (r *Runtime) sessionList(c echo.Context) error {
	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	sessions, err := common.SessionList(c.Request().Context(), r, user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
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

	return c.JSON(http.StatusOK, payload)
}

func (r *Runtime) sessionDelete(c echo.Context) error {
	sid := c.QueryParam("sid")
	if sid == "" {
		var p SidRequest
		err := c.Bind(&p)
		if err == nil && p.Sid != "" {
			sid = p.Sid
		}
	}
	if sid == "" {
		return c.String(http.StatusBadRequest, "missing sid")
	}

	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	err := common.SessionDel(c.Request().Context(), r, user, sid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

type stateRequest struct {
	Sid string `json:"sid"`
	Key string `json:"key"`
	Val any    `json:"val"`
}

func (r *Runtime) sessionStateGet(c echo.Context) error {
	var p stateRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	val, err := common.SessionStateGet(c.Request().Context(), r, user, p.Sid, p.Key)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"val": val})
}

func (r *Runtime) sessionStatePut(c echo.Context) error {
	var p stateRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	err = common.SessionStatePut(c.Request().Context(), r, user, p.Sid, p.Key, p.Val)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (r *Runtime) sessionStateDel(c echo.Context) error {
	var p stateRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	err = common.SessionStateDel(c.Request().Context(), r, user, p.Sid, p.Key)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
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

	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	spliced, err := common.SessionSplice(c.Request().Context(), r, user, p.Sid, p.Pos, p.Count, p.Fill)
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

	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	prompt, err := common.SessionPromptRender(c.Request().Context(), r, user, p.Sid, p.Agent)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"prompt": prompt})
}
