package api

import (
	"fmt"
	"net/http"

	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
	"github.com/hofstadter-io/hof/lib/consts"
	"github.com/labstack/echo/v4"
)

type fsPayload struct {
	Uri     string `json:"uri"`            // The OCI/Veg URI
	Path    string `json:"path,omitempty"` // Optional path override
	Sid     string `json:"sid,omitempty"`  // Session ID for context
	User    string `json:"user,omitempty"` // User for authorization
	Diff    bool   `json:"diff"`           // Whether to apply diff-based logic
	DiffUri string `json:"diffUri"`        // Base URI for diffing
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
	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = p.User
	}
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}
	fmt.Printf("fsStat: %s %s (user: %s)\n", p.Uri, p.Path, user)

	stat, err := common.FilesysStat(c.Request().Context(), r, user, p.Uri, p.Path, p.Sid, p.Diff)
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
	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = p.User
	}
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}

	content, err := common.FilesysRead(c.Request().Context(), r, user, p.Uri, p.Path, p.Sid, p.Diff)
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
	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = p.User
	}
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}
	fmt.Printf("fsList: %s %s (user: %s)\n", p.Uri, p.Path, user)

	entries, err := common.FilesysList(c.Request().Context(), r, user, p.Uri, p.Path, p.Sid, p.Diff)
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
	User    string `json:"user,omitempty"`
	Content string `json:"content"`
}

func (r *Runtime) fsWrite(c echo.Context) error {
	var p fsWriteRequest
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
	fmt.Printf("fsWrite: %s %s (user: %s)\n", p.Uri, p.Path, user)

	// p.Uri  ~ oci://host.docker.internal:5000/8c06ca51-3040-4d2a-951d-3c80e6ff084c%3A2?path%3D.%252Fextensions%252Fvscode%252Fextension%252Fsrc%252Fservices%252FfilesystemProvider.ts
	// p.Path ~ ./extensions/vscode/extension/src/services/filesystemProvider.ts

	nextUri, err := common.FilesysWrite(c.Request().Context(), r, user, p.Uri, p.Path, p.Content, p.Sid)
	if err != nil {
		fmt.Println("fsWrite.error:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsDeleteRequest struct {
	Uri  string `json:"uri"`
	Path string `json:"path"`
	Sid  string `json:"sid,omitempty"`
	User string `json:"user,omitempty"`
}

func (r *Runtime) fsDelete(c echo.Context) error {
	var p fsDeleteRequest
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

	nextUri, err := common.FilesysDelete(c.Request().Context(), r, user, p.Uri, p.Path, p.Sid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsMkdirRequest struct {
	Uri  string `json:"uri"`
	Path string `json:"path"`
	Sid  string `json:"sid,omitempty"`
	User string `json:"user,omitempty"`
}

func (r *Runtime) fsMkdir(c echo.Context) error {
	var p fsMkdirRequest
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

	nextUri, err := common.FilesysMkdir(c.Request().Context(), r, user, p.Uri, p.Path, p.Sid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsRenameRequest struct {
	Uri  string `json:"uri"`
	Src  string `json:"src"`
	Dst  string `json:"dst"`
	Sid  string `json:"sid,omitempty"`
	User string `json:"user,omitempty"`
}

func (r *Runtime) fsRename(c echo.Context) error {
	var p fsRenameRequest
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

	nextUri, err := common.FilesysRename(c.Request().Context(), r, user, p.Uri, p.Src, p.Dst, p.Sid)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"envUri": nextUri})
}

type fsCopyRequest struct {
	Uri  string `json:"uri"`
	Src  string `json:"src"`
	Dst  string `json:"dst"`
	Sid  string `json:"sid,omitempty"`
	User string `json:"user,omitempty"`
}

func (r *Runtime) fsCopy(c echo.Context) error {
	var p fsCopyRequest
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

	nextUri, err := common.FilesysCopy(c.Request().Context(), r, user, p.Uri, p.Src, p.Dst, p.Sid)
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
	user := c.Request().Header.Get(consts.VEG_USER_HEADER)
	if user == "" {
		user = p.User
	}
	if user == "" {
		user = consts.VEG_DEFAULT_USER
	}
	fmt.Printf("fsDiff: %v (user: %s)\n", p, user)

	diff, err := common.FilesysDiff(c.Request().Context(), r, user, p.DiffUri, p.Uri, p.Sid)
	if err != nil {
		fmt.Println("fs.Diff.error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	return c.JSON(http.StatusOK, diff)
}
