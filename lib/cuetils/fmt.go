package cuetils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/encoding/yaml"
)

func FormatOutput(val cue.Value, format string) (string, error) {
	switch format {
	case "cue", "CUE":
		return formatCue(val)

	case "json":
		return formatJson(val)

	case "yml", "yaml":
		return formatYaml(val)

	default:
		return "", fmt.Errorf("unknown output encoding %q", format)
	}

}

func formatCue(val cue.Value) (string, error) {
	syn := val.Syntax(
		cue.Final(),
		cue.ResolveReferences(true),
		cue.Concrete(true),
		cue.Definitions(false),
		cue.Hidden(false),
		cue.Optional(false),
		cue.Attributes(false),
		cue.Docs(false),
	)

	bs, err := format.Node(syn)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func formatJson(val cue.Value) (string, error) {
	var w bytes.Buffer
	d := json.NewEncoder(&w)
	d.SetIndent("", "  ")

	err := d.Encode(val)
	if _, ok := err.(*json.MarshalerError); ok {
		return "", err
	}

	return w.String(), nil
}

func formatYaml(val cue.Value) (string, error) {
	bs, err := yaml.Encode(val)

	if err != nil {
		return "", err
	}

	return string(bs), nil
}
