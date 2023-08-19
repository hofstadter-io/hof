package flags

type VetFlagpole struct {
	Expression []string
	Extensions bool
	List       bool
	Out        string
	Outfile    string
	Schema     string
	Concrete   bool
}

var VetFlags VetFlagpole
