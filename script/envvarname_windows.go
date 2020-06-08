package script

import "strings"

func envvarname(k string) string {
	return strings.ToLower(k)
}
