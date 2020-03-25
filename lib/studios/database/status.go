package database

import (
	"github.com/hofstadter-io/hof/pkg/util"
)

func Status() error {
	return util.SimpleGet("/studios/db/status")
}
