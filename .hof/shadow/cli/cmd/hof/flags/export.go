package flags

type ExportFlagpole struct {
	Expression []string
	List       bool
	Simplify   bool
	Out        string
	Outfile    string
	Escape     bool
	Comments   bool
}

var ExportFlags ExportFlagpole
