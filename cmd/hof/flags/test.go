package flags

type TestFlagpole struct {
	List        bool
	Keep        bool
	Suite       []string
	Tester      []string
	Environment []string
}

var TestFlags TestFlagpole
