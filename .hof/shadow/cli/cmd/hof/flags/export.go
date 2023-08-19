package flags

type ExportFlagpole struct {
	Expression []string
	List       bool
	Out        string
	Outfile    string
	Escape     bool
	Comments   bool
}

var ExportFlags ExportFlagpole
