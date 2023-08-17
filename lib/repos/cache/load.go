package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	// "strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/lib/repos/remote"
	"github.com/hofstadter-io/hof/lib/repos/utils"
)

const NoCacheEnvVar = "HOF_NOCACHE"
var noCache = false

const CustomCacheDirEnvVar = "HOF_CACHE"
var cacheBaseDir string
var modBaseDir string
var srcBaseDir string

// hacky, but we only want to sync once per repo per process
var syncedRepos *sync.Map
var cacheLock sync.Mutex

// CacheDir/{mod,src}

var debug = false

func init() {
	syncedRepos = new(sync.Map)
	nc := os.Getenv(NoCacheEnvVar)
	if nc != "" {
		noCache = true
	}

	e := os.Getenv(CustomCacheDirEnvVar)
	if e != "" {
		cacheBaseDir = e
	} else {
		d, err := os.UserCacheDir()
		if err != nil {
			return
		}

		// workaround for running in TestScript tool
		//if strings.HasPrefix(d, "/no-home") {
		//  d = strings.TrimPrefix(d, "/")
		//}

		// save to hof dir for cache across projects
		cacheBaseDir = filepath.Join(d, "hof")
	}
	modBaseDir = filepath.Join(cacheBaseDir, "mods")
	srcBaseDir = filepath.Join(cacheBaseDir, "src")
}

func SetCacheDir(basedir string) {
	cacheBaseDir = basedir
}

func Load(url, ver string) (_ billy.Filesystem, err error) {
	if debug {
		fmt.Println("cache.Load:", url, ver)
	}

	if !noCache {
		FS, err := Read(url, ver)
		if err == nil {
			// fmt.Println("cached:", url, ver)
			return FS, nil
		}
	}

	return Cache(url, ver)
}

func Read(url, ver string) (FS billy.Filesystem, err error) {
	if debug {
		fmt.Println("cache.Read:", url, ver)
	}
	remote, owner, repo := utils.ParseModURL(url)
	dir := SourceOutdir(remote, owner, repo)
	if ver != "" {
		dir = ModuleOutdir(remote, owner, repo, ver)
	}

	// check for existence
	_, err = os.Lstat(dir)
	if err != nil {
		return nil, err
	}

	// load into FS
	FS = osfs.New(dir)

	return FS, nil
}

func Cache(url, ver string) (billy.Filesystem, error) {
	if debug {
		fmt.Println("cache.Cache:", url, ver)
	}

	s, err := UpgradePseudoVersion(url, ver)
	if err != nil {
		return nil, err
	}
	ver = s

	if debug {
		fmt.Println("cache.Cache version resolve:", url, ver)
	}

	if ver == "" {
		rmt, err := remote.Parse(url)
		if err != nil {
			return nil, err
		}

		kind, err := rmt.Kind()
		if err != nil {
			return nil, err
		}

		switch kind {
		case remote.KindGit:
			return FetchRepoSource(url, ver)
		case remote.KindOCI:
			return FetchOCISource(url, ver)
		}
	}
	return CacheModule(url, ver)
}

func CacheModule(url, ver string) (billy.Filesystem, error) {
	if debug {
		fmt.Println("cache.CacheModule:", url, ver)
	}

	s, err := UpgradePseudoVersion(url, ver)
	if err != nil {
		return nil, err
	}
	ver = s

	if debug {
		fmt.Println("cache.CacheModule version resolve:", url, ver)
	}

	reg, owner, repo := utils.ParseModURL(url)
	dir := ModuleOutdir(reg, owner, repo, ver)

	// check for existing directory
	if _, err := os.Lstat(dir); err != nil {
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return nil, err
		}
		// doesnt' exist, so continue
	} else {
		// the directory exists
		if noCache {
			err = os.RemoveAll(dir)
			if err != nil {
				return nil, err
			}
		} else {
			// 
			return nil, fmt.Errorf("module already exists and noCache is not enabled, you should check for existence or use Load to be dynamic")
		}
	}

	rmt, err := remote.Parse(url)
	if err != nil {
		return nil, err
	}

	kind, err := rmt.Kind()
	if err != nil {
		return nil, err
	}

	switch kind {
	case remote.KindGit:
		cacheLock.Lock()
		defer cacheLock.Unlock()

		// we are smarter here and check to see if the tag already exists
		// this will both clone new & sync existing repos as needed
		// when ver != "", it will only fetch if the tag is not found
		_, err := FetchRepoSource(url, ver)
		if err != nil {
			return nil, err
		}

		// fmt.Println("MAKE COPY:", url, ver)
		ver, err = CopyRepoTag(url, ver)
		// fmt.Println("DONE MAKING:", url, ver)
		if err != nil {
			return nil, err
		}
	case remote.KindOCI:
		return FetchOCISource(url, ver)
	}


	return Read(url, ver)
}

