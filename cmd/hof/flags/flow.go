package flags

type FlowFlagpole struct {
	List       bool
	Docs       bool
	Flow       []string
	Tags       []string
	DebugTasks bool
}

var FlowFlags FlowFlagpole
