package flags

type EvalFlagpole struct {
	Expression    []string
	List          bool
	Out           string
	Outfile       string
	InlineImports bool
	Comments      bool
	All           bool
	Concrete      bool
	Attributes    bool
	Hidden        bool
	Optional      bool
}

var EvalFlags EvalFlagpole
