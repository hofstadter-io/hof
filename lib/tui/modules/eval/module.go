package eval

import "github.com/hofstadter-io/hof/lib/connector"

func New() connector.Connector {
	items := []any{
		NewEval(),
	}
	m := connector.New("Eval")
	m.Add(items)

	return m
}
