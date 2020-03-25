package database

import (
	"github.com/hofstadter-io/hof/pkg/util"
)

func Migrate() error {
	return util.SimpleGet("/studios/db/migrate")
}
