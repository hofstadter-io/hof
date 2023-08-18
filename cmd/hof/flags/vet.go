package flags

type VetFlagpole struct {
	Expression []string
	Extensions bool
	List       bool
	Out        string
	Schema     string
	Concrete   bool
}

var VetFlags VetFlagpole
