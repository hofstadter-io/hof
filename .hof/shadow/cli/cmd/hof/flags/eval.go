package flags

type EvalFlagpole struct {
	Expression    []string
	List          bool
	Simplify      bool
	Out           string
	Outfile       string
	InlineImports bool
	Comments      bool
	All           bool
	Concrete      bool
	Escape        bool
	Attributes    bool
	Definitions   bool
	Hidden        bool
	Optional      bool
	Resolve       bool
	Defaults      bool
	Final         bool
}

var EvalFlags EvalFlagpole
