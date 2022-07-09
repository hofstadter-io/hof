package flags

type FlowFlagpole struct {
	List     bool
	Docs     bool
	Flow     []string
	Tags     []string
	Progress bool
	Stats    bool
}

var FlowFlags FlowFlagpole