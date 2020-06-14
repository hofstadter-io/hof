package flags

type RootPflagpole struct {
	Labels             []string
	Config             string
	Secret             string
	ContextFile        string
	Context            string
	Global             bool
	Local              bool
	Input              []string
	InputFormat        string
	Output             []string
	OutputFormat       string
	Error              []string
	ErrorFormat        string
	Account            string
	Billing            string
	Project            string
	Workspace          string
	DatamodelDir       string
	ResourcesDir       string
	RuntimesDir        string
	Package            string
	Errors             bool
	Ignore             bool
	Simplify           bool
	Trace              bool
	Strict             bool
	Verbose            string
	Quiet              bool
	ImpersonateAccount string
	TraceToken         string
	LogHTTP            string
	RunWeb             bool
	RunTUI             bool
	RunREPL            bool
	Topic              string
	Example            string
	Tutorial           string
}

var RootPflags RootPflagpole
