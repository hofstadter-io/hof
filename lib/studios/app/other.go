package app

import (
	"github.com/hofstadter-io/hof/lib/util"
)

func Hello() error {
	return util.SimpleGet("/studios/app/hello")
}

func Generate() error {
	return util.SimpleGet("/studios/app/generate")
}

func Validate() error {
	return util.SimpleGet("/studios/app/validate")
}

func Versions() error {
	return util.SimpleGet("/studios/app/versions")
}
