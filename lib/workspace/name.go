package workspace

import (
	"fmt"
	"strings"
)

func CheckSplitModuleName(module string) (parts []string, err error) {
	fields := strings.Split(module, "/")

	if len(fields) != 3 {
		return nil, fmt.Errorf("Incorrect module format, should be 3 parts like github.com/hofstadter-io/hof")
	}

	// todo, check fields

	return fields, nil
}

