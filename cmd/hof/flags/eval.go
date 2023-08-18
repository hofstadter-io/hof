package flags

type EvalFlagpole struct {
	All        bool
	Concrete   bool
	Attributes bool
	Hidden     bool
	Optional   bool
}

var EvalFlags EvalFlagpole
