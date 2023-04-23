package remote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		desc    string
		mod     string
		outKind kind
	}{
		{
			desc:    "git github",
			mod:     "github.com/hofstadter-io/hof",
			outKind: kindGit,
		},
		{
			desc:    "git github private",
			mod:     "github.com/andrewhare/env",
			outKind: kindGit,
		},
		{
			desc:    "git not github",
			mod:     "git.kernel.org/pub/scm/bluetooth/bluez.git",
			outKind: kindGit,
		},
		{
			desc:    "oci",
			mod:     "gcr.io/distroless/static-debian11",
			outKind: kindOCI,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.desc, func(t *testing.T) {
			t.Parallel()

			out, err := Parse(c.mod)
			assert.NoError(t, err)
			assert.Equal(t, c.outKind, out.kind)
		})
	}
}
