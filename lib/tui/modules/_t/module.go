package {{ .name | lower }}

import "github.com/hofstadter-io/hof/lib/connector"

func New() connector.Connector {
	items := []any{
		New{{ .name | title }}(),
	}
	m := connector.New("{{ .name | title }}")
	m.Add(items)

	return m
}

