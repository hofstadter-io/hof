package runtime

import (
	"fmt"
	"maps"
	"net/http"
	"net/url"

	"google.golang.org/adk/session"

	"github.com/hofstadter-io/hof/lib/agent/agents"
	"github.com/hofstadter-io/hof/lib/agent/runtime/services/environ"
	"github.com/labstack/echo/v4"
)

type fsPayload struct {
	Uri     string `json:"uri"`
	Path    string `json:"path,omitempty"`
	Diff    bool   `json:"diff"`
	DiffUri string `json:"diffUri"`
	// User string `json:"user,omitempty"`
}

func fsOpen(c echo.Context) error {
	// fmt.Println("fsStat.called!")
	var opts environ.EnvironCreateOptions
	err := c.Bind(&opts)
	if err != nil {
		fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	envUri, err := environ.Client().Create(opts)
	if err != nil {
		fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	resp := map[string]any{
		"envUri": envUri,
	}

	return c.JSON(http.StatusOK, resp)
}

func fsStat(c echo.Context) error {
	var p fsPayload
	err := c.Bind(&p)
	if err != nil {
		fmt.Println("fsStat.bind.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	// fmt.Println("fsStat.called", p)

	u, err := url.Parse(p.Uri)
	if err != nil {
		fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	if p.Path == "" {
		p.Path = u.Query().Get("path")
	}

	stat, err := environ.Client().Stat(p.Uri, p.Path, p.Diff)
	if err != nil {
		// fmt.Println("fsStat.stat.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	// fmt.Println("fsStat.return", stat)

	return c.JSON(http.StatusOK, stat)
}

func (r *Runtime) fsRead(c echo.Context) error {
	var p fsPayload
	err := c.Bind(&p)
	if err != nil {
		fmt.Println("fsRead.bind.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	u, err := url.Parse(p.Uri)
	if err != nil {
		fmt.Println("fsRead.parse.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	if p.Path == "" {
		p.Path = u.Query().Get("path")
	}

	content, err := environ.Client().ReadFile(p.Uri, p.Path, p.Diff)
	if err != nil {
		// fmt.Println("fsRead.ReadFile.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	file := map[string]any{
		"uri":     p.Uri,
		"content": content,
	}

	return c.JSON(http.StatusOK, file)
}

func (r *Runtime) fsList(c echo.Context) error {
	// fmt.Println("fsList.called!")
	var p fsPayload
	err := c.Bind(&p)
	if err != nil {
		fmt.Println("fsList.bind.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	u, err := url.Parse(p.Uri)
	if err != nil {
		fmt.Println("fsList.parse.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	if p.Path == "" {
		p.Path = u.Query().Get("path")
	}

	entries, err := environ.Client().ReadDirectory(p.Uri, p.Path, p.Diff)
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, entries)
}

type fsWriteRequest struct {
	Uri     string `json:"uri"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (r *Runtime) fsWrite(c echo.Context) error {
	var p fsWriteRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	nextUri, err := environ.Client().WriteFile(p.Uri, p.Path, p.Content)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsDeleteRequest struct {
	Uri  string `json:"uri"`
	Path string `json:"path"`
}

func (r *Runtime) fsDelete(c echo.Context) error {
	var p fsDeleteRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	nextUri, err := environ.Client().Delete(p.Uri, p.Path, true)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsMkdirRequest struct {
	Uri  string `json:"uri"`
	Path string `json:"path"`
}

func (r *Runtime) fsMkdir(c echo.Context) error {
	var p fsMkdirRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	nextUri, err := environ.Client().CreateDirectory(p.Uri, p.Path)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsRenameRequest struct {
	Uri  string `json:"uri"`
	Src  string `json:"src"`
	Dst  string `json:"dst"`
}

func (r *Runtime) fsRename(c echo.Context) error {
	var p fsRenameRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	nextUri, err := environ.Client().Move(p.Uri, p.Src, p.Dst, true)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsCopyRequest struct {
	Uri  string `json:"uri"`
	Src  string `json:"src"`
	Dst  string `json:"dst"`
}

func (r *Runtime) fsCopy(c echo.Context) error {
	var p fsCopyRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	nextUri, err := environ.Client().Copy(p.Uri, p.Src, p.Dst, true)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsDiffRequest struct {
	PrevUri string `json:"prev"`
	NextUri string `json:"next"`
}

func (r *Runtime) fsDiff(c echo.Context) error {
	var p fsPayload
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	fmt.Println("fsDiff:", p)

	diff, err := environ.Client().DiffDirectory(p.DiffUri, p.Uri)
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	return c.JSON(http.StatusOK, diff)
}

func (r *Runtime) envList(c echo.Context) error {
	envs, err := environ.Client().ListEnvirons()
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, envs)
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
		AppName:   "veg",
		UserID:    "tony",
		SessionID: p.Sid,
	}
	sresp, err := r.S.Get(c.Request().Context(), sreq)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	sess := sresp.Session

	// 2. Get Agent Config
	agentName := p.Agent
	if agentName == "" {
		// Try to get agent from state or use a default if available
		// For now, if empty, we might need it passed or found in state
		v, _ := sess.State().Get("agent")
		if v != nil {
			agentName = v.(string)
		}
	}
	if agentName == "" {
		agentName = "veggie" // fallback default
	}

	agt, ok := r.Agentic.Agents[agentName]
	if !ok {
		return c.String(http.StatusNotFound, "agent not found: "+agentName)
	}

	// 3. Prepare State
	// TODO: We may need to walk the events backwards and process state changes inversely,
	// from the current state for the session (I don't think it's recorded, only the latest and delta,
	// we can leave this for later, leave a comment where it should go for now
	st := maps.Collect(sess.State().All())

	// 4. Render
	prompt, err := agents.RenderInstructionsWithNameAndState(r.Agentic, agt, agentName, st)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"prompt": prompt})
}
