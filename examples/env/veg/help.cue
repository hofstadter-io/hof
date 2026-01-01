@experiment(aliasv2)
package veg

import (
	"list"
	// "strings"
	"text/template"

	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

_l1: [...string]
_l1: [
	"flags",
	"src",
	"bins",
	"ctr",
	"fmtr",
	"svc",
	"dist",
	"info",
]
// _l2: [string]: [...string]
// _l2: { bins: ["multi"] }
// _l3: [string]: [string]: [...string]

infoData: {
	layout: {
		for _, l in _l1 {
			(l): list.Sort([for k, _ in root[l] {k}], {x: string, y: string, less: x < y})
		}
	}
	steps: list.Sort([for k, _ in env {k}], {x: string, y: string, less: x < y})
}

help: info.layout

info: {
	// example layout
	layout: template.Execute(layoutTmpl, infoData.layout)

	// list of steps
	steps: template.Execute(stepsTmpl, infoData.steps)

	// list of defs (?)
}

layoutTmpl: """
	{{ range $top, $items := . }}
	{{ $top }}
	{{- range $items }}
	  {{ . }}{{ end }}
	{{ end}}
	
	"""

stepsTmpl: """
	{{ range $top, $step := . }}
	{{ $step }}
	{{- end}}
	
	"""
