package gen

import (
	"bytes"
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
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
			// only want real errors
			// todo, we may want to inform the user better here
			if err != nil {
				if _, ok := err.(*hfmt.NoFormatterError); !ok {
					return err
				}
			}
		} else {
			fmtd, err = hfmt.FormatSource(F.Filepath, fmtd, "", nil, !F.FormattingDisabled)
			// only want real errors
			// todo, we may want to inform the user better here
			if err != nil {
				if _, ok := err.(*hfmt.NoFormatterError); !ok {
					return err
				}
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
		return FormatJson(val)

	case "yml", "yaml":
		return FormatYaml(val)

	case "xml":
		return FormatXml(val)

	case "toml":
		return FormatToml(val)

	default:
		return nil, fmt.Errorf("unknown output encoding %q", format)
	}
}

func (F *File) formatCue(val cue.Value) ([]byte, error) {

	// v := val

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

	return FormatCue(val, opts, F.Package)
}

func FormatCue(val cue.Value, opts []cue.Option, pkg string) ([]byte, error) {
	syn := val.Syntax(opts...)
	if pkg != "" {
		pkgDecl := &ast.Package {
			Name: ast.NewIdent(pkg),
		}
		decls := []ast.Decl{pkgDecl}
		// this could cause an issue?
		switch t := syn.(type) {
		case *ast.File:
			t.Decls = append(decls, t.Decls...)

		case *ast.StructLit:
			decls = append(decls, t.Elts...)
			f := &ast.File{
				Decls: decls,
			}
			syn = f
		case *ast.ListLit:
			decls = append(decls, t)
			f := &ast.File{
				Decls: decls,
			}
			syn = f
		}
	}

	bs, err := format.Node(syn)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func FormatJson(val cue.Value) ([]byte, error) {
	var w bytes.Buffer
	d := json.NewEncoder(&w)
	d.SetIndent("", "  ")

	err := d.Encode(val)
	if _, ok := err.(*json.MarshalerError); ok {
		return nil, err
	}

	return w.Bytes(), nil
}

func FormatYaml(val cue.Value) ([]byte, error) {
	bs, err := yaml.Encode(val)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func FormatToml(val cue.Value) ([]byte, error) {
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

func FormatXml(val cue.Value) ([]byte, error) {
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
