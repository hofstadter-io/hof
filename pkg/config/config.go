package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var c Config

type Config struct {
	CurrentContext string `yaml: "CurrentContext"`

	Contexts map[string]Context

}

type Context struct {
	APIKey   string `yaml: "APIKey"`
	Account  string `yaml: "Account"`
	Host     string `yaml: "Host"`
}

func Get() Config {
	load()
	return c
}

func GetConfigContext(context string) (Context, error) {
	load()

	ctx, ok := c.Contexts[context]
	if !ok {
		return Context{}, errors.New("Unknown Context: " + context)
	}

	// TODO remove once the old domain is retired
	if ctx.Host == "https://studios.studios.live.hofstadter.io" {
		ctx.Host = "https://studios.hofstadter.io"
		c.Contexts[context] = ctx
		write()
	}

	return ctx, nil
}

func GetCurrentContext() (Context) {
	load()

	context := c.CurrentContext
	ctx := c.Contexts[context]

	// TODO remove once the old domain is retired
	if ctx.Host == "https://studios.studios.live.hofstadter.io" {
		ctx.Host = "https://studios.hofstadter.io"
		c.Contexts[context] = ctx
		write()
	}

	return ctx
}

func load() error {
	err := viper.Unmarshal(&c)
	if err != nil {
		return err
	}

	return nil
}

func write() (err error) {

	B, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "error executing template\n")
	}

	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
			if home == "" {
					home = os.Getenv("USERPROFILE")
			}
	}

  dir := filepath.Join(home, ".hof")

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return errors.Wrap(err, "error writing hof config file\n")
	}

	err = ioutil.WriteFile(filepath.Join(dir, "hof.yaml"), B, 0644)
	if err != nil {
		return errors.Wrap(err, "error writing hof config file\n")
	}

	return nil
}
