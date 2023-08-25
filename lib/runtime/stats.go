package runtime

import (
	"bytes"
	"fmt"
	"time"
)

type RuntimeStats map[string]time.Duration

func (S RuntimeStats) Add(name string, dur time.Duration) {
	S[name] = dur
}
func (S RuntimeStats) String() string {
	var b bytes.Buffer

	order := []string{
		"cue/load",
		"data/load",
		"gen/load",
		"gen/run",
		"enrich/data",
		"enrich/gen",
		// "enrich/flow",
	}

	for _, o := range order {
		d, _ := S[o]
		fmt.Fprintf(&b, "%-16s%v\n", o, d.Round(time.Millisecond))
	}
	return b.String()
}

