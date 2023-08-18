package flags

type EvalFlagpole struct {
	Expression []string
	Extensions bool
	All        bool
	Concrete   bool
	Attributes bool
	Hidden     bool
	Optional   bool
}

var EvalFlags EvalFlagpole
