package flags

type VetFlagpole struct {
	Expression []string
	List       bool
	Out        string
	Outfile    string
	Concrete   bool
}

var VetFlags VetFlagpole
