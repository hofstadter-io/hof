package csp

import "cuelang.org/go/cue"

type Msg struct {
  Key string
  Val cue.Value
}
