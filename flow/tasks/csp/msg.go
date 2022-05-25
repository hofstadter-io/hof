package csp

import "cuelang.org/go/cue"

type Msg struct {
	Key string    `json:"key"`
	Val cue.Value `json:"val"`
}
