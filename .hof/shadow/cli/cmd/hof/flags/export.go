package flags

type ExportFlagpole struct {
	Expression []string
	Extensions bool
	List       bool
	Out        string
	Schema     string
	Escape     bool
}

var ExportFlags ExportFlagpole
