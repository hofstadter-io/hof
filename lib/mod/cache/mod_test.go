package cache

import (
	"fmt"
	"testing"
)

func TestSplitMod(t *testing.T) {
	tests := map[string]struct {
		mod      string
		expected string
	}{
		"simple":             {mod: "github.com/owner/repo", expected: "github.com/owner/repo"},
		"submodule":          {mod: "github.com/owner/repo/submodule", expected: "github.com/owner/repo"},
		"complex":            {mod: "gitlab.com/owner/repo.git/submodule", expected: "gitlab.com/owner/repo"},
		"subgroup":           {mod: "gitlab.com/owner/subgroup/repo.git", expected: "gitlab.com/owner/subgroup/repo"},
		"subgroup+submodule": {mod: "gitlab.com/owner/subgroup/repo.git/submodule", expected: "gitlab.com/owner/subgroup/repo"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rm, o, rp := splitMod(tc.mod)
			if fmt.Sprintf("%s/%s/%s", rm, o, rp) != tc.expected {
				t.Fatalf("expected: %v, got: %s/%s/%s", tc.expected, rm, o, rp)
			}
		})
	}
}
