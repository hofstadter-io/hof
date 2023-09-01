package help

import "github.com/hofstadter-io/hof/lib/connector"

func New() connector.Connector {
	items := []interface{}{
		NewHelp(),
	}
	m := connector.New("help")
	m.Add(items)

	return m
}
