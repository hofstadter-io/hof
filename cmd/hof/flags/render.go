package flags

type RenderFlagpole struct {
	Template []string
	Partial  []string
	Diff3    bool
}

var RenderFlags RenderFlagpole
