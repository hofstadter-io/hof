package flags

type GenFlagpole struct {
	Stats       bool
	Generator   []string
	Template    []string
	Partial     []string
	Diff3       bool
	NoFormat    bool
	KeepDeleted bool
	Watch       bool
	WatchFull   []string
	WatchFast   []string
	AsModule    string
	Outdir      string
}

var GenFlags GenFlagpole
