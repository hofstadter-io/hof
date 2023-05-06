package oci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirExcluded(t *testing.T) {
	cases := []struct {
		desc string
		d    Dir
		rels map[string]bool
	}{
		{
			desc: "simple",
			d: NewDir("", "/", []string{
				"foo",
				"/bar/baz",
			}),
			rels: map[string]bool{
				"foo":      true,
				"111":      false,
				"/bar/baz": true,
			},
		},
		{
			desc: "only permit specific files",
			d: NewDir("", "cue.mod", []string{
				"*",
				"!module.cue",
				"!sums.cue",
			}),
			rels: map[string]bool{
				"module.cue": false,
				"sums.cue":   false,
				"111":        true,
				"/bar/baz":   true,
			},
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.desc, func(t *testing.T) {
			for rel, expected := range c.rels {
				assert.Equal(t, expected, c.d.Excluded(rel), rel)
			}
		})
	}
}
