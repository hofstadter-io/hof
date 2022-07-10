package flags

type GenFlagpole struct {
	List         bool
	Stats        bool
	Generator    []string
	Template     []string
	Partial      []string
	Diff3        bool
	Watch        bool
	WatchGlobs   []string
	WatchXcue    []string
	CreateModule string
}

var GenFlags GenFlagpole
