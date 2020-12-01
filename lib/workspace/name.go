package workspace

import (
	"fmt"
	"strings"
)

func CheckSplitModuleName(module string) (parts []string, err error) {
	fields := strings.Split(module, "/")

	if len(fields) < 2 {
		return nil, fmt.Errorf("Incorrect module format, should be atleast 2 parts like domain.com/name. 3 parts is more common")
	}

	// todo, check fields

	return fields, nil
}

