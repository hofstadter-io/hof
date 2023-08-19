package flags

type DefFlagpole struct {
	Expression    []string
	Extensions    bool
	List          bool
	Out           string
	Outfile       string
	Schema        string
	InlineImports bool
	Attributes    bool
}

var DefFlags DefFlagpole
