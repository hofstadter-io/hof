package gen

import (
	"bytes"
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
	// "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/encoding/yaml"
	"github.com/clbanning/mxj"
	"github.com/BurntSushi/toml"

	hfmt "github.com/hofstadter-io/hof/lib/fmt"
)

type FmtConfig struct {
	Formatter string
	Config    interface{}
}

func (F *File) FormatRendered() (err error) {

	if !F.FormattingDisabled {
		// current content
		fmtd := F.RenderContent

		// check for custom config
		if F.FormattingConfig != nil && F.FormattingConfig.Config != nil {
			fmtd, err = hfmt.FormatSource(
				F.Filepath, fmtd,
				F.FormattingConfig.Formatter,
				F.FormattingConfig.Config,
				!F.FormattingDisabled,
			)
			if err != nil {
				return err
			}
		} else {
			fmtd, err = hfmt.FormatSource(F.Filepath, fmtd, "", nil, !F.FormattingDisabled)
			if err != nil {
				return err
			}
		}

		// update content
		F.RenderContent = fmtd
	}
	return nil
}

func (F *File) formatData(val cue.Value, format string) ([]byte, error) {
	switch format {
	case "cue":
		return F.formatCue(val)

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

func (F *File) formatCue(val cue.Value) ([]byte, error) {

	// v := val
	if F.Package != "" {
		fmt.Printf("File has package: %# v\n", F.Hidden)
		/*
		pkgDecl := &ast.Package {
			Name: ast.NewIdent(F.Package),
		}
		*/

	}

	opts := []cue.Option{
		cue.Concrete(F.Concrete),
		cue.Definitions(F.Definitions),
		cue.Optional(F.Optional),
		cue.Hidden(F.Hidden),
		cue.Attributes(F.Attributes),
		cue.Docs(F.Docs),
		// cue.InlineImports(F.InlineImports),
		cue.ErrorsAsValues(F.ErrorsAsValues),
	}
	if F.Final {
		opts = append(opts, cue.Final())
	}
	if F.Raw {
		opts = append(opts, cue.Raw())
	}

	syn := val.Syntax(opts...)

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

	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
