package remote

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		desc string
		mod  string
		out  Remote
	}{
		{
			desc: "git",
			mod:  "github.com/hofstadter-io/hof",
			out:  gitRemote{},
		},
		{
			desc: "oci",
			mod:  "gcr.io/distroless/static-debian11",
			out:  ociRemote{},
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.desc, func(t *testing.T) {
			out, err := Parse(c.mod)
			assert.NoError(t, err)
			assert.Equal(t, reflect.TypeOf(c.out), reflect.TypeOf(out))
		})
	}
}
