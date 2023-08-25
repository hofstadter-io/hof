package flags

type ChatFlagpole struct {
	Model       string
	System      []string
	Messages    []string
	Examples    []string
	Outfile     string
	Choices     int
	MaxTokens   int
	Temperature float64
	TopP        float64
	TopK        int
	Stop        []string
}

var ChatFlags ChatFlagpole
