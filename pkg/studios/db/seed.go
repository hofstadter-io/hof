package db

import (
	"github.com/hofstadter-io/hof/lib/util"
)

func Seed() error {
	return util.SimpleGet("/studios/db/seed")
}

