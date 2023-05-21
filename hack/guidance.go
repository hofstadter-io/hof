package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

var tmpl = `
{{ .foo }}

{{ gen "tell me about something" }}

{{ .cow }}
`

var d = map[string]any{
	"foo": "bar",
	"cow": "moo",
}

func main() {
	
	var buf bytes.Buffer

	var fs = template.FuncMap{
		"gen": func(prompt string) string {
			curr := buf.String()
			fmt.Println("GEN:", curr, prompt)
			return "called OpenAI..."
		},
	}


	t := template.Must(template.New("guidance").Funcs(fs).Parse(tmpl))

	err := t.Execute(&buf, d)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(buf.String())
}
