package runtime

import (
	"strconv"
)

// status checks the exit or status code from the last exec or http call
func (ts *Script) CmdStatus(neg int, args []string) {
	if len(args) != 1 {
		ts.Fatalf("usage: status <int>")
	}

	// Don't care
	if neg < 0 {
		return
	}

	// Check arg
	code, err := strconv.Atoi(args[0])
	if err != nil {
		ts.Fatalf("error: %v\nusage: stdin <int>", err)
	}

	// wanted different but got samd
	if neg > 0 && ts.status == code {
		ts.Fatalf("unexpected status match: %d", code)
	}

	if neg == 0 && ts.status != code {
		ts.Fatalf("unexpected status mismatch:  wated: %d  got %d", code, ts.status)
	}

}

