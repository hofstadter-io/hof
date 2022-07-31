package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	gofmt "go/format"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/encoding/yaml"
	"github.com/clbanning/mxj"
	"github.com/naoina/toml"

	hfmt "github.com/hofstadter-io/hof/lib/fmt"
)

var FORMAT_DISABLED = false

func init() {
	val := os.Getenv("HOF_FORMAT_DISABLED")
	if val == "true" || val == "1" {
		FORMAT_DISABLED=true
	}

	// gracefully init images / containers
	hfmt.GracefulInit()
}

func (F *File) FormatRendered() error {

	// If Golang only
	if strings.HasSuffix(F.Filepath, ".go") {
		fmtd, err := gofmt.Source(F.RenderContent)
		if err != nil {
			return err
		}

		F.RenderContent = fmtd
	}

	// hook into prettier / black ...
	// via container / system
	if !FORMAT_DISABLED {
		// inspect file settings to see if there is fmtr config...
		// try using hof/fmt containers, this is auto inference
		fmtd, err := hfmt.FormatSource(F.Filepath, F.RenderContent, "", nil)
		if err != nil {
			return err
		}

		F.RenderContent = fmtd
	}

	return nil
}

func formatData(val cue.Value, format string) ([]byte, error) {
	switch format {
	case "cue":
		return formatCue(val)

	case "json":
		return formatJson(val)

	case "yml", "yaml":
		return formatYaml(val)

	case "xml":
		return formatXml(val)

	case "toml":
		return formatToml(val)

	default:
		return nil, fmt.Errorf("unknown output encoding %q", format)
	}
}

func formatCue(val cue.Value) ([]byte, error) {

	syn := val.Syntax(
		cue.Final(),
		cue.Definitions(true),
		cue.Hidden(true),
		cue.Optional(true),
		cue.Attributes(true),
		cue.Docs(true),
	)

	bs, err := format.Node(syn)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func formatJson(val cue.Value) ([]byte, error) {
	var w bytes.Buffer
	d := json.NewEncoder(&w)
	d.SetIndent("", "  ")

	err := d.Encode(val)
	if _, ok := err.(*json.MarshalerError); ok {
		return nil, err
	}

	return w.Bytes(), nil
}

func formatYaml(val cue.Value) ([]byte, error) {
	bs, err := yaml.Encode(val)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func formatToml(val cue.Value) ([]byte, error) {
	v := make(map[string]interface{})
	err := val.Decode(&v)
	if err != nil {
		return nil, err
	}

	bs, err := toml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func formatXml(val cue.Value) ([]byte, error) {
	v := make(map[string]interface{})
	err := val.Decode(&v)
	if err != nil {
		return nil, err
	}

	mv := mxj.Map(v)
	bs, err := mv.XmlIndent("", "  ")
	if err != nil {
		return nil, err
	}
	return bs, nil
}
