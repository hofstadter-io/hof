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
	Definitions   bool
	Hidden        bool
	Optional      bool
	Final         bool
}

var EvalFlags EvalFlagpole
