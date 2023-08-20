package flags

type VetFlagpole struct {
	Expression  []string
	List        bool
	Out         string
	Outfile     string
	Concrete    bool
	Comments    bool
	Attributes  bool
	Definitions bool
	Hidden      bool
	Optional    bool
}

var VetFlags VetFlagpole
