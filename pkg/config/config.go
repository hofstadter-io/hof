package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const CONFIG_FILENAME = "hof-mod.yaml"

type Config struct {
	Name string
	Type string
	Path string

	Entrypoint string

	Deps []Dependency
}

type Dependency struct {
	Package string
	Version string
	Replace string
}

func LoadModuleConfig(dir string) (*Config, error) {
	filename := filepath.Join(dir, CONFIG_FILENAME)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
