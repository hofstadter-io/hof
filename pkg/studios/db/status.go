package db

import (
	"github.com/hofstadter-io/hof/lib/util"
)

func Status() error {
	return util.SimpleGet("/studios/db/status")
}
