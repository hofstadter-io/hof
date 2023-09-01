package eval

import "github.com/hofstadter-io/hof/lib/connector"

func New() connector.Connector {
	items := []interface{}{
		NewEvalPage(),
	}

	m := connector.New("eval")
	m.Add(items)

	return m
}
