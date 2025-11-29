package runtime

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
	"github.com/labstack/echo/v4"
	"google.golang.org/adk/session"
)

type fileStat struct {
	Sid  string `json:"sid"`
	Path string `json:"path,omitempty"`

	Hash  string `json:"hash,omitempty"`
	Dir   bool   `json:"dir"`
	Size  int    `json:"size"`
	Ctime int64  `json:"ctime"`
	Mtime int64  `json:"mtime"`
}

type fsPayload struct {
	Sid  string `json:"sid"`
	Pos  string `json:"pos,omitempty"`
	Path string `json:"path,omitempty"`
	User string `json:"user,omitempty"`
}

func (r *Runtime) fsCommon(sid string) (session.Session, *dagger.Directory, error) {
	resp, err := r.S.Get(r.Ctx, &session.GetRequest{
		AppName:   r.AppName,
		UserID:    "tony",
		SessionID: sid,
	})
	if err != nil {
		return resp.Session, nil, err
	}
	dagId, _ := resp.Session.State().Get("dagger")
	dagDir := r.Dagger.LoadDirectoryFromID(dagger.DirectoryID(dagId.(string)))

	return resp.Session, dagDir, nil
}

func (r *Runtime) fsStat(c echo.Context) error {
	// fmt.Println("fsStat.called!")
	var p fsPayload
	err := c.Bind(&p)
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	// fmt.Println("fsStat.payload", p)

	session, dagDir, err := r.fsCommon(p.Sid)
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	// does it exist?
	ok, err := dagDir.Exists(r.Ctx, p.Path, dagger.DirectoryExistsOpts{})
	if err != nil || !ok {
		// fmt.Println("error/ok:", err, !ok)
		return c.String(http.StatusBadRequest, "bad request")
	}

	// get first event time
	var t int64
	if session.Events().Len() > 0 {
		t = session.Events().At(0).Timestamp.UnixMilli()
	}

	stat := fileStat{
		Sid:   p.Sid,
		Path:  p.Path,
		Ctime: t,
		Mtime: t,
	}

	// is it a directory?
	ok, _ = dagDir.Exists(r.Ctx, p.Path, dagger.DirectoryExistsOpts{ExpectedType: dagger.ExistsTypeDirectoryType})
	if ok {
		stat.Dir = true
	} else {
		stat.Size, err = dagDir.File(p.Path).Size(r.Ctx)
		if err != nil {
			// fmt.Println("error:", err)
			return c.String(http.StatusBadRequest, "bad request")
		}
	}

	return c.JSON(http.StatusOK, stat)
}

func (r *Runtime) fsRead(c echo.Context) error {
	// fmt.Println("fsRead.called!")
	var p fsPayload
	err := c.Bind(&p)
	if err != nil {
		fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}
	// fmt.Println("fsRead.payload", p)

	session, dagDir, err := r.fsCommon(p.Sid)
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	ok, err := dagDir.Exists(r.Ctx, p.Path, dagger.DirectoryExistsOpts{ExpectedType: dagger.ExistsTypeRegularType})
	if err != nil || !ok {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	contents, err := dagDir.File(p.Path).Contents(r.Ctx)
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	// could possibly do better mtime by walking events, but it should be gettable from dagger one would think
	// get first event time
	var t int64
	if session.Events().Len() > 0 {
		t = session.Events().At(0).Timestamp.UnixMilli()
	}

	file := map[string]any{
		"sid":      p.Sid,
		"path":     p.Path,
		"contents": contents,
		"ctime":    t,
		"mtime":    t,
	}

	return c.JSON(http.StatusOK, file)
}

func (r *Runtime) fsList(c echo.Context) error {
	// fmt.Println("fsList.called!")
	var p fsPayload
	err := c.Bind(&p)
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	// fmt.Println("fsList.sid/path", p.Sid, p.Path)

	session, dagDir, err := r.fsCommon(p.Sid)
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	ok, err := dagDir.Exists(r.Ctx, p.Path, dagger.DirectoryExistsOpts{ExpectedType: dagger.ExistsTypeDirectoryType})
	if err != nil || !ok {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	entries, err := dagDir.Directory(p.Path).Entries(r.Ctx)
	if err != nil {
		// fmt.Println("error:", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	pairs := []any{}
	for _, e := range entries {
		ok, _ := dagDir.Exists(r.Ctx, filepath.Join(p.Path, e), dagger.DirectoryExistsOpts{ExpectedType: dagger.ExistsTypeDirectoryType})
		// if !strings.HasPrefix(e, "/") {
		// 	e = "/" + e
		// }
		pairs = append(pairs, []any{e, ok})
	}

	// could possibly do better mtime by walking events, but it should be gettable from dagger one would think
	// get first event time
	var t int64
	if session.Events().Len() > 0 {
		t = session.Events().At(0).Timestamp.UnixMilli()
	}

	file := map[string]any{
		"sid":     p.Sid,
		"path":    p.Path,
		"entries": pairs,
		"ctime":   t,
		"mtime":   t,
	}

	// fmt.Println("fsList.pairs", p)
	return c.JSON(http.StatusOK, file)
}

func (r *Runtime) fsWrite(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (r *Runtime) fsDelete(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// func (r *Runtime) fsDiff(c echo.Context) error {
// 	var p fsPayload
// 	err := c.Bind(&p)
// 	if err != nil {
// 		return c.String(http.StatusBadRequest, "bad request")
// 	}

// 	sess, dagDir, err := r.fsCommon(p.Sid)
// 	if err != nil {
// 		fmt.Println("error:", err)
// 		return c.String(http.StatusBadRequest, "bad request")
// 	}
// 	origId, _ := sess.State().Get("origfs")
// 	origDir := dag.LoadDirectoryFromID(dagger.DirectoryID(origId.(string)))

// 	return c.String(http.StatusOK, "Hello, World!")
// }

// func (r *Runtime) fsExport(c echo.Context) error {
// 	var p fsPayload
// 	err := c.Bind(&p)
// 	if err != nil {
// 		return c.String(http.StatusBadRequest, "bad request")
// 	}

// 	sess, dagDir, err := r.fsCommon(p.Sid)
// 	if err != nil {
// 		fmt.Println("error:", err)
// 		return c.String(http.StatusBadRequest, "bad request")
// 	}
// 	origId, _ := sess.State().Get("origfs")
// 	origDir := dag.LoadDirectoryFromID(dagger.DirectoryID(origId.(string)))

// 	// get diff, then export

// 	// dagDir.Diff()

// 	return c.String(http.StatusOK, "Hello, World!")
// }

func DiffDirectories(ctx context.Context, prev, next *dagger.Directory) (map[string]any, error) {
	changes := next.Changes(prev)
	// fmt.Println("session.diff.debug", origId, dagId, maps.Collect(resp.Session.State().All()))

	addpaths, err := changes.AddedPaths(ctx)
	if err != nil {
		return nil, err
	}

	modpaths, err := changes.ModifiedPaths(ctx)
	if err != nil {
		return nil, err
	}

	delpaths, err := changes.RemovedPaths(ctx)
	if err != nil {
		return nil, err
	}

	// ensure absolute
	for i, fp := range addpaths {
		if !strings.HasPrefix(fp, "/") {
			fp = "/" + fp
			addpaths[i] = fp
		}
	}

	for i, fp := range modpaths {
		if !strings.HasPrefix(fp, "/") {
			fp = "/" + fp
			modpaths[i] = fp
		}
	}

	for i, fp := range delpaths {
		if !strings.HasPrefix(fp, "/") {
			fp = "/" + fp
			delpaths[i] = fp
		}
	}

	// gather current files
	diffFiles := make(map[string]string)
	diffErrs := make([]string, 0)
	for _, fp := range addpaths {
		// it appears this is how they are distinguished in the string, which is better than nothing!
		if strings.HasSuffix(fp, "/") {
			continue
		}

		f, err := next.File(fp).Contents(ctx)
		if err != nil {
			err = fmt.Errorf("while reading %q: %w", fp, err)
			diffErrs = append(diffErrs, err.Error())
		}

		diffFiles[fp] = f
	}
	for _, fp := range modpaths {
		// it appears this is how they are distinguished in the string, which is better than nothing!
		if strings.HasSuffix(fp, "/") {
			continue
		}

		f, err := next.File(fp).Contents(ctx)
		if err != nil {
			err = fmt.Errorf("while reading %q: %w", fp, err)
			diffErrs = append(diffErrs, err.Error())
		}

		diffFiles[fp] = f
	}

	pfile := changes.AsPatch()
	patch, err := pfile.Contents(ctx)

	return map[string]any{
		"addpaths": addpaths,
		"modpaths": modpaths,
		"delpaths": delpaths,
		"patch":    patch,
		"files":    diffFiles,
		"diffErrs": diffErrs,
	}, nil
}
