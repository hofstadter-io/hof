package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var ChatFlagSet *pflag.FlagSet

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

func SetupChatFlags(fset *pflag.FlagSet, fpole *ChatFlagpole) {
	// flags

	fset.StringVarP(&(fpole.Model), "model", "M", "gpt-3.5-turbo", "LLM model to use [gpt-3.5-turbo,gpt-4,bard,chat-bison]")
	fset.StringArrayVarP(&(fpole.System), "system", "s", nil, "string or path to the system prompt for the LLM, concatenated")
	fset.StringArrayVarP(&(fpole.Messages), "message", "m", nil, "string or path to a message for the LLM")
	fset.StringArrayVarP(&(fpole.Examples), "example", "e", nil, "string or path to an example pair for the LLM")
	fset.StringVarP(&(fpole.Outfile), "outfile", "O", "", "path to write the output to")
	fset.IntVarP(&(fpole.Choices), "choices", "N", 1, "param: choices or N (openai)")
	fset.IntVarP(&(fpole.MaxTokens), "max-tokens", "", 256, "param: MaxTokens")
	fset.Float64VarP(&(fpole.Temperature), "temp", "", 0.8, "param: temperature")
	fset.Float64VarP(&(fpole.TopP), "topp", "", 0.42, "param: TopP")
	fset.IntVarP(&(fpole.TopK), "topk", "", 40, "param: TopK (google)")
	fset.StringArrayVarP(&(fpole.Stop), "stop", "", nil, "param: Stop (openai)")
}

func init() {
	ChatFlagSet = pflag.NewFlagSet("Chat", pflag.ContinueOnError)

	SetupChatFlags(ChatFlagSet, &ChatFlags)

}
