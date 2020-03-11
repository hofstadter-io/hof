package io

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/hofstadter-io/hof-lang/parser"
	"github.com/naoina/toml"
	"gopkg.in/yaml.v2"
)

/*
Where's your docs doc?!
*/
func InferDataContentType(data []byte) (contentType string, err error) {

	// TODO: look for unique symbols in the data
	// but always try to unmarshal to be sure

	var obj interface{}

	err = json.Unmarshal(data, &obj)
	if err == nil {
		return "json", nil
	}

	err = yaml.Unmarshal(data, &obj)
	if err == nil {
		return "yaml", nil
	}

	err = xml.Unmarshal(data, &obj)
	if err == nil {
		return "yaml", nil
	}

	err = toml.Unmarshal(data, &obj)
	if err == nil {
		return "toml", nil
	}

	_, err = parser.ParseReader("", bytes.NewReader(data))
	if err == nil {
		return "hof", nil
	}

	return "", errors.New("[IDCT] unknown content type")

	return
}

/*
Where's your docs doc?!
*/
func InferFileContentType(filename string) (contentType string, err error) {

	// assume files have correct extensions
	// TODO use 'filepath.Ext()'
	ext := filepath.Ext(filename)[1:]
	switch ext {

	case "json":
		return "json", nil

	case "toml":
		return "toml", nil

	case "yaml", "yml":
		return "yaml", nil

	case "xml":
		return "xml", nil

	case "hof":
		return "hof", nil

	default:
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", err
		}
		return InferDataContentType(data)
	}

	return
}

// HOFSTADTER_BELOW
