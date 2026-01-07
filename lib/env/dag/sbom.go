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

type hashSBOMConfig struct {
	Kind   string    `json:"$kind"`
	Name   string    `json:"name"`
	Path   string    `json:"path"`
	Format string    `json:"format"`
	Data   cue.Value `json:"data"`
}

type hashSBOMIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashSBOMConfig
	file *dagger.File
}

func (idx *hashSBOMIndex) Key() string {
	if idx.cfg == nil {
		return "#sbom.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#sbom.%s.%s", idx.cfg.Kind, mk)
	}
	return fmt.Sprintf("#sbom.%s.%s", idx.cfg.Kind, idx.cfg.Path)
}

func (d *Dag) HashCuefigSBOM(step cue.Value) (*dagger.File, string, error) {
	d.mx.RLock()
	var cfg hashSBOMConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, "", fmt.Errorf("while decoding hashCuefigSBOM: %w", err)
	}

	// index for query and create if not found
	idx := &hashSBOMIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashSBOMIndex)
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

func (d *Dag) HashDaggerSBOM(step cue.Value) (*dagger.File, string, error) {
	d.mx.RLock()
	var cfg hashSBOMConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, "", fmt.Errorf("while decoding hashDaggerSBOM: %w", err)
	}

	// index for query and create if not found
	idx := &hashSBOMIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashSBOMIndex)
		return ix.file, ix.cfg.Path, nil
	}

	// 1. Get the actual dagger object from the CUE value
	// We use the existing dispatchers to get the *dagger.Type
	var obj json.Marshaler
	var dErr error

	k := cfg.Data.LookupPath(cue.ParsePath("$kind"))
	if !k.Exists() {
		return nil, "", fmt.Errorf("missing $kind in DaggerSBOM data: %v", cfg.Data)
	}
	ks, _ := k.String()

	switch ks {
	case "#container", "#hostImage", "#dockerBuild":
		obj, dErr = d.Container(cfg.Data, false)
	case "#exportImageFile":
		obj, _, dErr = d.HashExportImageFile(cfg.Data)
	case "#exportImage":
		obj, _, dErr = d.HashExportImage(cfg.Data)
	case "#publishImage":
		obj, _, dErr = d.HashPublishImage(cfg.Data)

	case "#service":
		obj, _, dErr = d.Service(cfg.Data, false)
	case "#hostService":
		obj, _, dErr = d.HashHostService(cfg.Data)
	case "#hostTunnel":
		obj, _, dErr = d.HashHostTunnel(cfg.Data)
	case "#hostSocket":
		obj, _, dErr = d.HashHostSocket(cfg.Data)

	case "#file", "#hostFile":
		obj, _, dErr = d.File(cfg.Data, false)
	case "#exportFile":
		obj, _, dErr = d.HashExportFile(cfg.Data)

	case "#dir", "#hostDir", "#gitRepo":
		obj, _, dErr = d.Dir(cfg.Data, false)
	case "#exportDir":
		obj, _, dErr = d.HashExportDir(cfg.Data)

	case "#secret":
		obj, dErr = d.hashSecret(cfg.Data)
	case "#cache":
		obj, dErr = d.hashCache(cfg.Data)

	case "#cuefigSBOM":
		obj, _, dErr = d.HashCuefigSBOM(cfg.Data)
	case "#daggerSBOM":
		obj, _, dErr = d.HashDaggerSBOM(cfg.Data)

	default:
		return nil, "", fmt.Errorf("unsupported $kind %q in DaggerSBOM data", ks)
	}

	if dErr != nil {
		return nil, "", fmt.Errorf("while resolving Dagger object for SBOM: %w", dErr)
	}

	// 2. Marshal to JSON
	jsonBS, err := obj.MarshalJSON()
	if err != nil {
		return nil, "", fmt.Errorf("while marshaling Dagger object to JSON: %w", err)
	}

	// 3. Handle formats
	var bs []byte
	var cerr error

	switch cfg.Format {
	case "json":
		bs = jsonBS

	case "yaml":
		var tmp any
		cerr = json.Unmarshal(jsonBS, &tmp)
		if cerr == nil {
			bs, cerr = yaml.Marshal(tmp)
		}

	case "toml":
		var tmp any
		cerr = json.Unmarshal(jsonBS, &tmp)
		if cerr == nil {
			bs, cerr = toml.Marshal(tmp)
		}

	case "cue":
		// This might be tricky if we want "nice" CUE, but for now we can just use the JSON
		bs = jsonBS

	default:
		return nil, "", fmt.Errorf("unsupported format %q for hashDaggerSBOM", cfg.Format)
	}

	if cerr != nil {
		return nil, "", fmt.Errorf("while transcoding Dagger SBOM to %q: %w", cfg.Format, cerr)
	}

	f := d.dag.File(cfg.Path, string(bs))

	// memoize
	idx.file = f
	d.cat.Store(idx, idx)

	return idx.file, idx.cfg.Path, nil
}
