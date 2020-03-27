package gen

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

type GeneratorStats struct {
	NumNew       int
	NumSame      int
	NumSkipped   int
	NumDeleted   int
	NumWritten   int
	NumErr       int
	TotalFiles   int

	NumModified       int
	NumModifiedRender int
	NumModifiedOutput int
	NumModifiedDiff3  int
	NumConflicted     int

	CueLoadingTime time.Duration
	RenderingTime  time.Duration
	TotalTime      time.Duration
}

type FileStats struct {
	// using 0 (false) and 1 (true) for easier summation code below
	IsNew      int
	IsSame     int
	IsSkipped  int
	IsWritten  int
	IsErr      int

	IsModified       int
	IsModifiedRender int
	IsModifiedOutput int
	IsModifiedDiff3  int
	IsConflicted     int

	RenderingTime  time.Duration
	CompareTime    time.Duration
	TotalTime      time.Duration
}

func (S *GeneratorStats) CalcTotals(G *Generator) error {
	// Start with own fields
	var sum time.Time
	sum = sum.Add(S.CueLoadingTime)
	sum = sum.Add(S.RenderingTime)

	S.TotalTime = sum.Sub(time.Time{})
	S.TotalFiles = len(G.Files)

	// Sum across files
	for _, file := range G.Files {
		S.NumNew += file.IsNew
		S.NumSame += file.IsSame
		S.NumSkipped += file.IsSkipped
		S.NumWritten += file.IsWritten
		S.NumErr += file.IsErr

		S.NumModified += file.IsModified
		S.NumModifiedRender += file.IsModifiedRender
		S.NumModifiedOutput += file.IsModifiedOutput
		S.NumModifiedDiff3 += file.IsModifiedDiff3
		S.NumConflicted += file.IsConflicted
	}

	return nil
}

func (S *GeneratorStats) String() string {
	var b bytes.Buffer
	var err error

	// Parse Template
	t := template.Must(template.New("stats").Parse(statsTemplate))

	// Round timings
	S.CueLoadingTime = S.CueLoadingTime.Round(time.Microsecond)
	S.RenderingTime  = S.RenderingTime.Round(time.Microsecond)
	S.TotalTime      = S.TotalTime.Round(time.Microsecond)

	// Render template
	err = t.Execute(&b, S)
	if err != nil {
		return fmt.Sprint(err)
	}

	return b.String()
}

const statsTemplate = `
NumNew              {{ .NumNew }}
NumSame             {{ .NumSame }}
NumSkipped          {{ .NumSkipped }}
NumDeleted          {{ .NumDeleted }}
NumWritten          {{ .NumWritten }}
NumErr              {{ .NumErr }}

NumModified         {{ .NumModified }}
NumModifiedRender   {{ .NumModifiedRender }}
NumModifiedOutput   {{ .NumModifiedOutput }}
NumModifiedDiff3    {{ .NumModifiedDiff3 }}
NumConflicted       {{ .NumConflicted }}

CueLoadingTime      {{ .CueLoadingTime }}
RenderingTime       {{ .RenderingTime }}
TotalTime           {{ .TotalTime }}
`
