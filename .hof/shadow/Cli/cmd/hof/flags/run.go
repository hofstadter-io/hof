package flags

type RunFlagpole struct {
	List        bool
	Info        bool
	Suite       []string
	Runner      []string
	Environment []string
	Data        []string
	Workdir     string
}

var RunFlags RunFlagpole
