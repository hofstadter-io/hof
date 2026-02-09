package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
	"github.com/hofstadter-io/hof/lib/consts"
	"github.com/labstack/echo/v4"
)

type fsPayload struct {
	Uri     string `json:"uri"`
	Path    string `json:"path,omitempty"`
	Sid     string `json:"sid,omitempty"`
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

	envUri, err := environ.Client().Create(&opts)
	if err != nil {
		fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	resp := map[string]any{
		"envUri": envUri,
	}

	return c.JSON(http.StatusOK, resp)
}

func (r *Runtime) fsStat(c echo.Context) error {
	var p fsPayload
	err := c.Bind(&p)
	if err != nil {
		fmt.Println("fsStat.bind.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	fmt.Printf("fsStat: %s %s\n", p.Uri, p.Path)

	u, err := url.Parse(p.Uri)
	if err != nil {
		fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	if p.Path == "" {
		p.Path = u.Query().Get("path")
	}

	stat, err := common.FilesysStat(c.Request().Context(), r, consts.VEG_DEFAULT_USER, p.Uri, p.Path, p.Sid, p.Diff)
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
	fmt.Printf("fsRead: %s %s\n", p.Uri, p.Path)

	u, err := url.Parse(p.Uri)
	if err != nil {
		fmt.Println("fsRead.parse.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	if p.Path == "" {
		p.Path = u.Query().Get("path")
	}

	content, err := common.FilesysRead(c.Request().Context(), r, consts.VEG_DEFAULT_USER, p.Uri, p.Path, p.Sid, p.Diff)
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
	var p fsPayload
	err := c.Bind(&p)
	if err != nil {
		fmt.Println("fsList.bind.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	fmt.Printf("fsList: %s %s\n", p.Uri, p.Path)

	u, err := url.Parse(p.Uri)
	if err != nil {
		fmt.Println("fsList.parse.error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	if p.Path == "" {
		p.Path = u.Query().Get("path")
	}

	entries, err := common.FilesysList(c.Request().Context(), r, consts.VEG_DEFAULT_USER, p.Uri, p.Path, p.Sid, p.Diff)
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, entries)
}

type fsWriteRequest struct {
	Uri     string `json:"uri"`
	Path    string `json:"path"`
	Sid     string `json:"sid,omitempty"`
	Content string `json:"content"`
}

func (r *Runtime) fsWrite(c echo.Context) error {
	var p fsWriteRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	fmt.Printf("fsWrite: %s %s\n", p.Uri, p.Path)

	nextUri, err := common.FilesysWrite(c.Request().Context(), r, consts.VEG_DEFAULT_USER, p.Uri, p.Path, p.Content, p.Sid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsDeleteRequest struct {
	Uri  string `json:"uri"`
	Path string `json:"path"`
	Sid  string `json:"sid,omitempty"`
}

func (r *Runtime) fsDelete(c echo.Context) error {
	var p fsDeleteRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	nextUri, err := common.FilesysDelete(c.Request().Context(), r, consts.VEG_DEFAULT_USER, p.Uri, p.Path, p.Sid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsMkdirRequest struct {
	Uri  string `json:"uri"`
	Path string `json:"path"`
	Sid  string `json:"sid,omitempty"`
}

func (r *Runtime) fsMkdir(c echo.Context) error {
	var p fsMkdirRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	nextUri, err := common.FilesysMkdir(c.Request().Context(), r, consts.VEG_DEFAULT_USER, p.Uri, p.Path, p.Sid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsRenameRequest struct {
	Uri string `json:"uri"`
	Src string `json:"src"`
	Dst string `json:"dst"`
	Sid string `json:"sid,omitempty"`
}

func (r *Runtime) fsRename(c echo.Context) error {
	var p fsRenameRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	nextUri, err := common.FilesysRename(c.Request().Context(), r, consts.VEG_DEFAULT_USER, p.Uri, p.Src, p.Dst, p.Sid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsCopyRequest struct {
	Uri string `json:"uri"`
	Src string `json:"src"`
	Dst string `json:"dst"`
	Sid string `json:"sid,omitempty"`
}

func (r *Runtime) fsCopy(c echo.Context) error {
	var p fsCopyRequest
	err := c.Bind(&p)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	nextUri, err := common.FilesysCopy(c.Request().Context(), r, consts.VEG_DEFAULT_USER, p.Uri, p.Src, p.Dst, p.Sid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsDiffRequest struct {
	PrevUri string `json:"prev"`
	NextUri string `json:"next"`
}

func fsDiff(c echo.Context) error {
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
