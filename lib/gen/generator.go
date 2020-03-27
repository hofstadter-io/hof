package gen

import (
	"fmt"

	"cuelang.org/go/cue"
)

// A map to generators
type Generators map[string]*Generator

// A generator pulled from the cue instances
type Generator struct {
	// Label in Cuelang
	Name string

	// These will be set externally
	In map[string]interface{}
	Out []map[string]interface{}

	// Files and the shadow dir for doing neat things
	Files map[string]*File
	Shadow map[string]*File

	// Status for this generator and processing
	Stats *GeneratorStats

	// Cuelang related, also set externally
	CueValue         cue.Value
}

func NewGenerator(label string, value cue.Value) *Generator{
	return &Generator {
		Name: label,
		CueValue: value,
		Files: make(map[string]*File),
		Shadow: make(map[string]*File),
		Stats: &GeneratorStats{},
	}
}

func (G *Generator) GenerateFiles() error {
	errs := []error{}

	// Todo, make this a parallel work queue
	for _, F := range G.Files {
		F.ShadowFile = G.Shadow[F.Filename]
		err := F.Render()
		if err != nil {
			errs = append(errs, fmt.Errorf("In file %q, error %w", F.Filename, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Errors while rendering files:\n%v\n", errs)
	}

	return nil
}

