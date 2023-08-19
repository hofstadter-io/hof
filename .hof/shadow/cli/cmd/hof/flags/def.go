package flags

type DefFlagpole struct {
	Expression    []string
	List          bool
	Out           string
	Outfile       string
	InlineImports bool
	Comments      bool
	Attributes    bool
}

var DefFlags DefFlagpole
