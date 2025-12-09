package stores

import "dagger.io/dagger"

// Abstract fileSys/execEnv backed by Dagger as a unified interface
// for both REST and websocket handlers
// This serves two purposes
// 1. where we proxy/translate between vscode and dagger (standalone without ADK, do we care about separation at this point, trend is towards putting more hof/veg more UI in vscode)
// 2. what we attach to an Agent Session / Event
// Implementations need to store Hash->ID after changes in mutating functions
// ... if you want to save space and such in your Session state
// ... we should probably do something similar for files? (perhaps what artifacts are for in the first place?)
// maybe some of these need to return that, but then we are getting close to Daggerland
// We also need to be mindful where/how we mirror changes to the session event history when appropriate
type Container interface {
	// Dagger refs
	ID() (string, error) // hash ref to a dagger ID stored else where
	Load(id string) (*dagger.Container, error)

	// Exec related
	Exec(args ...string) (ExecResult, error)

	// VS Code
	Stat(path string) (FileStat, error)
	ReadFile(path string) (string, error)
	WriteFile(path, content string) error

	CreateDirectory(path string) error
	ReadDirectory(path string) ([]Dirent, error)

	Rename(src, dst string) error
	Copy(src, dst string) error
	Delete(path string, recursive bool) error

	Watch(path string, recursive bool)
}

type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type Dirent struct {
	Path string
	Dir  bool
}

type FileStat struct {
	Dir   bool
	Ctime int
	Mtime int
	Size  int
}
