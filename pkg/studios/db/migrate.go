package db

import (
	"github.com/hofstadter-io/hof/lib/util"
)

func Migrate() error {
	return util.SimpleGet("/studios/db/migrate")
}
