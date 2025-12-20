package runtime

import (
	"fmt"
	"net/http"
	"net/url"

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
