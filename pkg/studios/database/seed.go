package database

import (
	"github.com/hofstadter-io/hof/pkg/util"
)

func Seed() error {
	return util.SimpleGet("/studios/db/seed")
}

