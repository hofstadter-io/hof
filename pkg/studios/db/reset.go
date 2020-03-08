package db

import (
	"github.com/hofstadter-io/hof/lib/util"
)

func Reset(hard bool) error {
	if hard {
		return util.SimpleGet("/studios/db/reset?hard=yes")
	}

	return util.SimpleGet("/studios/db/reset")
}
