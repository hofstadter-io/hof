package dag

import (
	"fmt"

	"cuelang.org/go/cue"
)

type hashCmdConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`

	Tasks  map[string]cue.Value `json:"tasks"`
	Hooks  map[string]cue.Value `json:"hooks"`
	Config map[string]any       `json:"config"`
}

func (d *Dag) DecodeHashCmd(step cue.Value) (*hashCmdConfig, error) {
	var cfg hashCmdConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashCmd: %w", err)
	}

	return &cfg, nil
}

type hashTaskConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`

	Steps  [][]cue.Value        `json:"steps"`
	Hooks  map[string]cue.Value `json:"hooks"`
	Config map[string]any       `json:"config"`
}

func (d *Dag) DecodeHashTask(step cue.Value) (*hashTaskConfig, error) {
	var cfg hashTaskConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashTask: %w", err)
	}

	return &cfg, nil
}
