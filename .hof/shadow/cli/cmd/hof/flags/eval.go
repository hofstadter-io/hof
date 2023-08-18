package flags

type EvalFlagpole struct {
	Expression []string
	Extensions bool
	List       bool
	Out        string
	Schema     string
	All        bool
	Concrete   bool
	Attributes bool
	Hidden     bool
	Optional   bool
}

var EvalFlags EvalFlagpole
