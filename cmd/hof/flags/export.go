package flags

type ExportFlagpole struct {
	Expression []string
	Extensions bool
	Escape     bool
	Out        string
}

var ExportFlags ExportFlagpole
