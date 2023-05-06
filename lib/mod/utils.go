package mod

import (
	"fmt"
	"strings"
)

func ValidateModURL(mod string) error {
	parts := strings.Split(mod, "/")	
	if len(parts) < 2 {
		return fmt.Errorf("error: modules require one or more '/', you provided %q", mod)
	}
	if !strings.Contains(parts[0], ".") {
		return fmt.Errorf("error: the first part of a module path must be a domain, you provided %q", mod)
	}

	return nil
}
