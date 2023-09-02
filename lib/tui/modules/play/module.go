package play

import "github.com/hofstadter-io/hof/lib/connector"

func New() connector.Connector {
	items := []interface{}{
		NewPlay(),
	}
	m := connector.New("Play")
	m.Add(items)

	return m
}

