package home

import "github.com/hofstadter-io/hof/lib/connector"

func New() connector.Connector {
	items := []interface{}{
		NewHome(),
	}
	m := connector.New("home")
	m.Add(items)

	return m
}
