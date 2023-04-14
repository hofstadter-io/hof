package runtime

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

type RuntimeStats struct {
	CueLoadingTime time.Duration
	GenLoadingTime time.Duration
	GenRunningTime time.Duration
}

func (S *RuntimeStats) String() string {
	var b bytes.Buffer
	var err error

	// Parse Template
	t := template.Must(template.New("stats").Parse(runtimeStatsTemplate))

	// Round timings
	S.CueLoadingTime = S.CueLoadingTime.Round(time.Microsecond)
	S.GenLoadingTime = S.GenLoadingTime.Round(time.Microsecond)
	S.GenRunningTime = S.GenRunningTime.Round(time.Microsecond)

	// Render template
	err = t.Execute(&b, S)
	if err != nil {
		return fmt.Sprint(err)
	}

	return b.String()
}

const runtimeStatsTemplate = `
CueLoadingTime      {{ .CueLoadingTime }}
GenLoadingTime      {{ .GenLoadingTime }}
GenRunningTime      {{ .GenRunningTime }}
`

