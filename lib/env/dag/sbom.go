package dag

import (
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/naoina/toml"
	"gopkg.in/yaml.v3"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/env"
)

type hashCuefigSBOMConfig struct {
	Kind   string    `json:"$kind"`
	Name   string    `json:"name"`
	Path   string    `json:"path"`
	Format string    `json:"format"`
	Data   cue.Value `json:"data"`
}

type hashCuefigSBOMIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashCuefigSBOMConfig
	file *dagger.File
}

func (idx *hashCuefigSBOMIndex) Key() string {
	if idx.cfg == nil {
		return "#cuefigSBOM.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#cuefigSBOM.%s.%s", idx.cfg.Kind, mk)
	}
	return fmt.Sprintf("#cuefigSBOM.%s.%s", idx.cfg.Kind, idx.cfg.Path)
}

func (d *Dag) HashCuefigSBOM(step cue.Value, noCache bool) (*dagger.File, string, error) {
	d.mx.RLock()
	var cfg hashCuefigSBOMConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, "", fmt.Errorf("while decoding HashCuefigSBOM: %w", err)
	}

	// index for query and create if not found
	idx := &hashCuefigSBOMIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashCuefigSBOMIndex)
		return ix.file, ix.cfg.Path, nil
	}

	var bs []byte
	var cerr error

	// TODO: implement CuefigSBOM generation
	switch cfg.Format {
	case "cue":
		s, cerr := cuetils.ValueToSyntaxString(
			cfg.Data,
			cue.Final(),
			// cue.Concrete(true),
			cue.Definitions(true),
			cue.Hidden(false),
			cue.Optional(false),
			cue.Attributes(true),
			// cue.Docs(false),
		)
		if cerr != nil {
			return nil, "", fmt.Errorf("while printing CUE for CuefigSBOM: %w", cerr)
		}
		bs = []byte(s)

	case "json":
		bs, cerr = json.MarshalIndent(cfg.Data, "", "  ")
		if cerr != nil {
			return nil, "", fmt.Errorf("while marshaling JSON for CuefigSBOM: %w", cerr)
		}

	case "yaml":
		bs, cerr = yaml.Marshal(cfg.Data)
		if cerr != nil {
			return nil, "", fmt.Errorf("while marshaling YAML for CuefigSBOM: %w", cerr)
		}

	case "toml":
		bs, cerr = toml.Marshal(cfg.Data)
		if cerr != nil {
			return nil, "", fmt.Errorf("while marshaling TOML for CuefigSBOM: %w", cerr)
		}

	default:
		return nil, "", fmt.Errorf("unsupported format %q for hashCuefigSBOM", cfg.Format)
	}

	f := d.dag.File(cfg.Path, string(bs))

	// memoize
	idx.file = f
	d.cat.Store(idx, idx)

	return idx.file, idx.cfg.Path, nil
}
