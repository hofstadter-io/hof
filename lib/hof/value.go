package hof

import (
	"cuelang.org/go/cue"
)

// wrapper around a cue.Value to create a shared singleton
type Value struct {
	root     cue.Value
	value    cue.Value
	path     string
}

func WrapValue(root cue.Value, path string) *Value {
	return &Value{
		root: root,
		value: root.LookupPath(cue.ParsePath(path)),
		path: path,
	}
}

func (v *Value) CueValue() cue.Value {
	return v.value
}

func (v *Value) RootValue() cue.Value {
	return v.root
}

func (v *Value) LookupPath(path string) *Value {
	if v.path != "" {
		path = v.path + "." + path
	}
	return &Value{
		value: v.value,
		path: path,
	}
}

func (v *Value) FillPath(path string, fill any) {
	if v.path != "" {
		path = v.path + "." + path
	}
	v.value = v.value.FillPath(cue.ParsePath(path), fill)
}

func (v *Value) Err() error {
	return v.CueValue().Err()
}

func (v *Value) Decode(in any) error {
	return v.CueValue().Decode(in)
}

func (v *Value) Path() string {
	return v.path
}

func (v *Value) Exists() bool {
	return v.CueValue().Exists()
}

func (v *Value) IsConcrete() bool {
	return v.CueValue().IsConcrete()
}

func (v *Value) List() (cue.Iterator, error) {
	return v.CueValue().List()
}

func (v *Value) Fields(opts ...cue.Option) (*cue.Iterator, error) {
	return v.CueValue().Fields(opts...)
}

func (v *Value) Context() *cue.Context {
	return v.value.Context()
}
