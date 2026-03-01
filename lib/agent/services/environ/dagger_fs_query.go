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

// we'll need a way to specify other Uri, and if not set, use this "genesis/0" assumption
func (le *localEnviron) getDagDir(envUri string, diff bool) (*dagger.Directory, error) {
	_, env, err := le.LookupEnviron(envUri)
	if err != nil {
		return nil, fmt.Errorf("while looking up environ(%s): %w", envUri, err)
	}
	d := env.Directory(".")
	if diff {
		origUri := ReplaceTag(envUri, "0")
		if envUri == origUri {
			return d, nil
		}
		_, orig, err := le.LookupEnviron(origUri)
		if err != nil {
			return nil, fmt.Errorf("while looking up orig environ(%s): %w", envUri, err)
		}
		od := orig.Directory(".")
		d = od.Diff(d)
	}

	return d, nil
}

func (le *localEnviron) Stat(envUri, path string, diff bool) (*FileStat, error) {
	// fmt.Println("le.Stat", envUri, path)
	d, err := le.getDagDir(envUri, diff)
	if err != nil {
		return nil, fmt.Errorf("while getting base fs(%s): %s %v %w", envUri, path, diff, err)
	}

	// more fucking reshaping... seriously, fuck vscode for having shitty Uri implementation
	parseUri := envUri
	if !strings.Contains(parseUri, "://") {
		parseUri = "oci://" + parseUri
	}
	ruri, err := url.Parse(parseUri)
	if err != nil {
		return nil, fmt.Errorf("while parsing uri(%s): %w", envUri, err)
	}

	if path == "" {
		path = ruri.Query().Get("path")
	}

	ok, err := d.Exists(le.ctx, path, dagger.DirectoryExistsOpts{})
	if err != nil || !ok {
		return nil, fmt.Errorf("while looking up path(%s): %v %w", path, ok, err)
	}
	if !ok {
		return nil, fmt.Errorf("file not found(%s): %v %w", path, ok, err)
	}

	stat := &FileStat{
		Uri:  envUri,
		Path: path,
		// Ctime: row.CreatedAt.UnixMicro(), // this should be the first time the file showed up
		// Mtime: row.UpdatedAt.UnixMicro(), // this should be the current env time (create/update likely always the same, unless we add names or allow changing uri?)
	}
	// fmt.Println("le.Stat.found", row, stat)

	// is it a directory?
	ok, _ = d.Exists(le.ctx, path, dagger.DirectoryExistsOpts{ExpectedType: dagger.ExistsTypeDirectoryType})
	if ok {
		stat.Dir = true
	} else {
		stat.Size, err = d.File(path).Size(le.ctx)
		if err != nil {
			// fmt.Println("error:", err)
			return nil, fmt.Errorf("while getting size for file(%s): %w", path, err)
		}
	}

	return stat, nil
}

// TODO, add offset/limit here
func (le *localEnviron) ReadFile(envUri, path string, diff bool) (string, error) {
	// fmt.Println("le.Stat", envUri, path)
	d, err := le.getDagDir(envUri, diff)
	if err != nil {
		return "", fmt.Errorf("while getting base fs(%s): %s %v %w", envUri, path, diff, err)
	}

	// more fukcing reshaping... seriously, fuck vscode for having shitty Uri implementation
	parseUri := envUri
	if !strings.Contains(parseUri, "://") {
		parseUri = "oci://" + parseUri
	}
	ruri, err := url.Parse(parseUri)
	if err != nil {
		return "", fmt.Errorf("while parsing uri(%s): %w", envUri, err)
	}
	if path == "" {
		path = ruri.Query().Get("path")
	}

	content, err := d.File(path).Contents(le.ctx, dagger.FileContentsOpts{
		OffsetLines: 0,
		LimitLines:  2000,
	})
	if err != nil {
		return "", fmt.Errorf("while getting contents(%s): %w", path, err)
	}

	return content, nil
}

func (le *localEnviron) ReadDirectory(envUri, path string, diff bool) (*DirList, error) {
	fmt.Println("le.ReadDirectory.input", envUri, path, diff)
	d, err := le.getDagDir(envUri, diff)
	if err != nil {
		return nil, fmt.Errorf("while getting base fs(%s): %s %v %w", envUri, path, diff, err)
	}
	// fmt.Println("le.ReadDirectory.lookup", table)

	// more fucking reshaping... seriously, fuck vscode for having shitty Uri implementation
	// we need to move this to vscode, it should not be handled in the environ service
	parseUri := envUri
	if !strings.Contains(parseUri, "://") {
		parseUri = "oci://" + parseUri
	}
	ruri, err := url.Parse(parseUri)
	if err != nil {
		return nil, fmt.Errorf("while parsing uri(%s): %w", envUri, err)
	}
	if path == "" {
		path = ruri.Query().Get("path")
	}
	fmt.Println("le.ReadDirectory.vars", ruri, path, diff)

	envEntries, err := d.Directory(path).Entries(le.ctx)
	if err != nil {
		return nil, fmt.Errorf("while listing directory(%s): %w", path, err)
	}
	fmt.Println("le.ReadDirectory.envEntries", envEntries)

	entries := []DirEntry{}
	for _, e := range envEntries {
		realPath := filepath.Join(path, e)
		// they must exist since we already go them
		ok, _ := d.Exists(le.ctx, realPath, dagger.DirectoryExistsOpts{ExpectedType: dagger.ExistsTypeDirectoryType})

		// Fallback: Exists(Directory) can return false for implicit directories in Diff layers
		// So we double check by trying to list it
		if !ok {
			_, err := d.Directory(realPath).Entries(le.ctx)
			if err == nil {
				ok = true
			}
		}

		fmt.Printf("le.ReadDirectory.entry: path=%q real=%q isDir=%v\n", path, realPath, ok)
		entries = append(entries, DirEntry{Name: e, Dir: ok})
	}

	return &DirList{
		Uri:     envUri,
		Path:    path,
		Entries: entries,
	}, nil
}

func (le *localEnviron) GrepDirectory(envUri, pattern string, diff bool) ([]dagger.SearchResult, error) {
	d, err := le.getDagDir(envUri, diff)
	if err != nil {
		return nil, fmt.Errorf("while getting base fs(%s): %s %v %w", envUri, pattern, diff, err)
	}

	results, err := d.Search(le.ctx, pattern, dagger.DirectorySearchOpts{
		Limit:       42,
		SkipIgnored: true,
	})
	if err != nil {
		return nil, fmt.Errorf("while searching environment(%s): %w", envUri, err)
	}

	return results, nil
}

func (le *localEnviron) GlobDirectory(envUri, pattern string, diff bool) ([]string, error) {
	d, err := le.getDagDir(envUri, diff)
	if err != nil {
		return nil, fmt.Errorf("while getting base fs(%s): %s %v %w", envUri, pattern, diff, err)
	}

	results, err := d.Glob(le.ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("while globbing environment(%s): %w", envUri, err)
	}
	if len(results) > 42 {
		results = append(results[:42], "... results truncated")
	}

	return results, nil
}

func (le *localEnviron) FindAgentFiles(envUri string) (map[string]string, error) {
	_, env, err := le.LookupEnviron(envUri)
	if err != nil {
		return nil, fmt.Errorf("while looking up environment(%s): %w", envUri, err)
	}

	results := make(map[string]string)
	agents, err := env.Directory(".").Glob(le.ctx, "**/AGENTS.md")
	if err != nil {
		return nil, fmt.Errorf("while globbing for agents (%s): %w", envUri, err)
	}
	for _, path := range agents {
		contents, err := env.Directory(".").File(path).Contents(le.ctx)
		if err != nil {
			return nil, fmt.Errorf("while reading agents.md (%s): %w", envUri, err)
		}
		results[path] = contents
	}

	claude, err := env.Directory(".").Glob(le.ctx, "**/CLAUDE.md")
	if err != nil {
		return nil, fmt.Errorf("while globbing for claude (%s): %w", envUri, err)
	}
	for _, path := range claude {
		contents, err := env.Directory(".").File(path).Contents(le.ctx)
		if err != nil {
			return nil, fmt.Errorf("while reading claude.md (%s): %w", envUri, err)
		}
		results[path] = contents
	}

	gemini, err := env.Directory(".").Glob(le.ctx, "**/GEMINI.md")
	if err != nil {
		return nil, fmt.Errorf("while globbing for gemini (%s): %w", envUri, err)
	}
	for _, path := range gemini {
		contents, err := env.Directory(".").File(path).Contents(le.ctx)
		if err != nil {
			return nil, fmt.Errorf("while reading gemini.md (%s): %w", envUri, err)
		}
		results[path] = contents
	}

	return results, nil
}

func (le *localEnviron) Watch(envUri, path string, excludes []string, recursive bool) (any, error) {

	return nil, nil
}

func (le *localEnviron) DiffDirectory(prevUri, nextUri string) (*DiffInfo, error) {
	fmt.Printf("le.DiffDirectory: prev=%s next=%s\n", prevUri, nextUri)
	if prevUri == "" {
		prevUri = ReplaceTag(nextUri, "0")
	}

	_, prev, err := le.LookupEnviron(prevUri)
	if err != nil {
		fmt.Println("dirdiff.lookup.prev.error:", err)
		return nil, fmt.Errorf("while looking up environment(%s): %w", prevUri, err)
	}

	prevDir := prev.Directory(".")
	_, next, err := le.LookupEnviron(nextUri)
	if err != nil {
		fmt.Println("dirdiff.lookup.next.error:", err)
		return nil, fmt.Errorf("while looking up environment(%s): %w", nextUri, err)
	}
	nextDir := next.Directory(".")

	fmt.Println("le.DiffDirectory: calculating changes", prevUri, nextUri)
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

	fmt.Printf("le.DiffDirectory: paths count: add=%d mod=%d del=%d\n", len(addpaths), len(modpaths), len(delpaths))

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

	fmt.Println("le.DiffDirectory: generating patch")
	pfile := changes.AsPatch()
	patch, err := pfile.Contents(le.ctx)
	if err != nil {
		return nil, fmt.Errorf("while getting patch contents(%s): %w", nextUri, err)
	}

	// DO THIS LAST, so we don't break other reads above

	// ensure absolute (should we be doing this?) forgot why we needed it in the first place... probably some other thing making it into the path in vscode or something
	// somewhere we lose the `/` on the beginning of an absolute path and that broke something
	// ~~Now turing this off because (1) source control explorer is borked (2) we see diffs outside of the workdir~~
	// workdir, err := next.Workdir(le.ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("while looking up workdir(%s): %w", nextUri, err)
	// }

	// Far too many different intertwining path knots, we need to untangle it and make a good pattern (<img>?path=...&sid=...) because we also need sid everywhere it's tied to an env
	// Apparently not having / in the front means we see the image:Npath/to/file (unable to separate)
	for i, fp := range addpaths {
		if !strings.HasPrefix(fp, "/") {
			fp = "/" + fp
			addpaths[i] = fp
		}
		// addpaths[i] = filepath.Join(workdir, fp)
	}

	for i, fp := range modpaths {
		if !strings.HasPrefix(fp, "/") {
			fp = "/" + fp
			modpaths[i] = fp
		}
		// modpaths[i] = filepath.Join(workdir, fp)
	}

	for i, fp := range delpaths {
		if !strings.HasPrefix(fp, "/") {
			fp = "/" + fp
			delpaths[i] = fp
		}
		// delpaths[i] = filepath.Join(workdir, fp)
	}
	if !strings.Contains(prevUri, "://") {
		prevUri = "oci://" + prevUri
	}
	if !strings.Contains(nextUri, "://") {
		nextUri = "oci://" + nextUri
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
