package flags

type TestFlagpole struct {
	List        bool
	Suite       []string
	Tester      []string
	Environment []string
}

var TestFlags TestFlagpole
