package flags

type EvalFlagpole struct {
	Expression    []string
	Extensions    bool
	List          bool
	Out           string
	Outfile       string
	Schema        string
	InlineImports bool
	Comments      bool
	All           bool
	Concrete      bool
	Attributes    bool
	Hidden        bool
	Optional      bool
}

var EvalFlags EvalFlagpole
