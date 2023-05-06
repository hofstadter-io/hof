package remote

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := []struct {
		desc    string
		mod     string
		outKind Kind
	}{
		{
			desc:    "git github",
			mod:     "github.com/hofstadter-io/hof",
			outKind: KindGit,
		},
		{
			desc:    "git github private",
			mod:     "github.com/andrewhare/env",
			outKind: KindGit,
		},
		{
			desc:    "git not github",
			mod:     "git.kernel.org/pub/scm/bluetooth/bluez.git",
			outKind: KindGit,
		},
		{
			desc:    "oci",
			mod:     "gcr.io/distroless/static-debian11",
			outKind: KindOCI,
		},
		{
			desc:    "oci",
			mod:     "us-central1-docker.pkg.dev/hof-io--develop/testing/test",
			outKind: KindOCI,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.desc, func(t *testing.T) {
			t.Parallel()

			out, err := Parse(c.mod)
			require.NoError(t, err)
			assert.Equal(t, c.outKind, out.kind)
		})
	}
}
