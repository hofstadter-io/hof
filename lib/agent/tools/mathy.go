package tools

import (
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type SumArgs struct {
	A int `json:"a"` // an integer to sum
	B int `json:"b"` // another integer to sum
}
type SumResult struct {
	Sum int `json:"sum"` // the sum of two integers
}

func NewSummer() (tool.Tool, error) {
	handler := func(ctx tool.Context, input SumArgs) (SumResult, error) {
		return SumResult{Sum: input.A + input.B}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "sum",
		Description: "sums two integers",
	}, handler)
}

type SubArgs struct {
	A int `json:"a"` // an integer to sum
	B int `json:"b"` // another integer to sum
}
type SubResult struct {
	Sum int `json:"sum"` // the sum of two integers
}

func NewSubber() (tool.Tool, error) {
	handler := func(ctx tool.Context, input SumArgs) (SumResult, error) {
		return SumResult{Sum: input.A - input.B}, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "sub",
		Description: "subtract two integers",
	}, handler)
}
