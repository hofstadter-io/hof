package cuetils

import (
	"cuelang.org/go/cue"
)

// TODO, improve merge strategy
func AttrToMap(A cue.Attribute) (m map[string]string) {
	m = make(map[string]string)
	for i := 0; i < A.NumArgs(); i++ {
		key, val := A.Arg(i)
		m[key] = val
	}
	return m
}
