package environ

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

type FileStat struct {
	Uri   string `json:"uri"`
	Path  string `json:"path"`
	Dir   bool   `json:"dir"`
	Size  int    `json:"size"`
	Hash  string `json:"hash,omitempty"`
	Ctime int64  `json:"ctime,omitempty"`
	Mtime int64  `json:"mtime,omitempty"`
}

type DirList struct {
	Uri     string     `json:"uri"`
	Path    string     `json:"path"`
	Entries []DirEntry `json:"entries"`
}

type DirEntry struct {
	Name string `json:"name"`
	Dir  bool   `json:"dir"`
}

type DiffInfo struct {
	PrevUri string `json:"prev"`
	NextUri string `json:"next"`

	AddPaths []string `json:"addPaths"`
	ModPaths []string `json:"modPaths"`
	DelPaths []string `json:"delPaths"`

	Patch string            `json:"patch"`
	Files map[string]string `json:"files"`

	Errors []string `json:"errors"`
}

func (le *localEnviron) Stat(envUri, path string) (*FileStat, error) {
	// fmt.Println("le.Stat", envUri, path)
	row, env, err := le.lookupEnviron(envUri)
	if err != nil {
		return nil, fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	ok, err := env.Exists(le.ctx, path, dagger.ContainerExistsOpts{})
	if err != nil || !ok {
		return nil, fmt.Errorf("while looking up path(%s): %v %w", path, ok, err)
	}

	ruri, err := url.Parse(envUri)
	if err != nil || !ok {
		return nil, fmt.Errorf("while parsing uri(%s): %w", envUri, err)
	}

	if path == "" {
		path = ruri.Query().Get("path")
	}

	stat := &FileStat{
		Uri:   envUri,
		Path:  path,
		Ctime: row.CreateAt.UnixMicro(), // this should be the first time the file showed up
		Mtime: row.UpdateAt.UnixMicro(), // this should be the current env time (create/update likely always the same, unless we add names or allow changing uri?)
	}
	// fmt.Println("le.Stat.found", row, stat)

	// is it a directory?
	ok, _ = env.Directory(".").Exists(le.ctx, path, dagger.DirectoryExistsOpts{ExpectedType: dagger.ExistsTypeDirectoryType})
	if ok {
		stat.Dir = true
	} else {
		stat.Size, err = env.File(path).Size(le.ctx)
		if err != nil {
			// fmt.Println("error:", err)
			return nil, fmt.Errorf("while getting size for file(%s): %w", path, err)
		}
	}

	return stat, nil
}

func (le *localEnviron) ReadFile(envUri, path string) (string, error) {

	_, env, err := le.lookupEnviron(envUri)
	if err != nil {
		return "", fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}
	// more fukcing reshaping... seriously, fuck vscode for having shitty Uri implementation
	ruri, err := url.Parse(envUri)
	if err != nil {
		return "", fmt.Errorf("while parsing uri(%s): %w", envUri, err)
	}
	if path == "" {
		path = ruri.Query().Get("path")
	}

	content, err := env.File(path).Contents(le.ctx, dagger.FileContentsOpts{})
	if err != nil {
		return "", fmt.Errorf("while getting contents(%s): %w", path, err)
	}

	return content, nil
}

func (le *localEnviron) ReadDirectory(envUri, path string) (*DirList, error) {
	// fmt.Println("le.ReadDirectory.input", envUri, path)
	_, env, err := le.lookupEnviron(envUri)
	if err != nil {
		return nil, fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	// more fukcing reshaping... seriously, fuck vscode for having shitty Uri implementation
	// we need to move this to vscode, it should not be handled in the environ service
	ruri, err := url.Parse(envUri)
	if err != nil {
		return nil, fmt.Errorf("while parsing uri(%s): %w", envUri, err)
	}
	if path == "" {
		path = ruri.Query().Get("path")
	}

	envEntries, err := env.Directory(path).Entries(le.ctx)
	if err != nil {
		return nil, fmt.Errorf("while listing directory(%s): %w", path, err)
	}
	// fmt.Println("le.ReadDirectory.envEntries", envEntries)

	entries := []DirEntry{}
	for _, e := range envEntries {
		realPath := filepath.Join(path, e)
		ok, _ := env.Directory(".").Exists(le.ctx, realPath, dagger.DirectoryExistsOpts{ExpectedType: dagger.ExistsTypeDirectoryType})
		// fmt.Println("le.ReadDirectory.entry", realPath, ok)
		entries = append(entries, DirEntry{Name: e, Dir: ok})
	}

	return &DirList{
		Uri:     envUri,
		Path:    path,
		Entries: entries,
	}, nil
}

func (le *localEnviron) GrepDirectory(envUri, pattern string) ([]dagger.SearchResult, error) {
	_, env, err := le.lookupEnviron(envUri)
	if err != nil {
		return nil, fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	results, err := env.Directory(".").Search(le.ctx, pattern, dagger.DirectorySearchOpts{
		Limit:       100,
		SkipIgnored: true,
		// Paths:       []string{workdir},
	})

	return results, err
}

func (le *localEnviron) Watch(envUri, path string, excludes []string, recursive bool) (any, error) {

	return nil, nil
}

func (le *localEnviron) DiffDirectory(prevUri, nextUri string) (*DiffInfo, error) {
	_, prev, err := le.lookupEnviron(prevUri)
	if err != nil {
		return nil, fmt.Errorf("while looking up environment(%s): %w", prevUri, err)
	}
	prevDir := prev.Directory(".")
	_, next, err := le.lookupEnviron(nextUri)
	if err != nil {
		return nil, fmt.Errorf("while looking up environment(%s): %w", nextUri, err)
	}
	nextDir := next.Directory(".")

	changes := nextDir.Changes(prevDir)

	addpaths, err := changes.AddedPaths(le.ctx)
	if err != nil {
		return nil, fmt.Errorf("while getting AddPaths(%s): %w", nextUri, err)
	}

	modpaths, err := changes.ModifiedPaths(le.ctx)
	if err != nil {
		return nil, fmt.Errorf("while getting ModPaths(%s): %w", nextUri, err)
	}

	delpaths, err := changes.RemovedPaths(le.ctx)
	if err != nil {
		return nil, fmt.Errorf("while getting DelPaths(%s): %w", nextUri, err)
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

		f, err := next.File(fp).Contents(le.ctx)
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

		f, err := next.File(fp).Contents(le.ctx)
		if err != nil {
			err = fmt.Errorf("while reading %q: %w", fp, err)
			diffErrs = append(diffErrs, err.Error())
		}

		diffFiles[fp] = f
	}

	pfile := changes.AsPatch()
	patch, err := pfile.Contents(le.ctx)
	if err != nil {
		return nil, fmt.Errorf("while getting patch contents(%s): %w", nextUri, err)
	}

	di := &DiffInfo{
		PrevUri:  prevUri,
		NextUri:  nextUri,
		AddPaths: addpaths,
		ModPaths: modpaths,
		DelPaths: delpaths,
		Patch:    patch,
		Files:    diffFiles,
		Errors:   diffErrs,
	}

	return di, nil
}
