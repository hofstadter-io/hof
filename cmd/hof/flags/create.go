package flags

type CreateFlagpole struct {
	Input     []string
	Generator []string
	Outdir    string
	Exec      bool
}

var CreateFlags CreateFlagpole
