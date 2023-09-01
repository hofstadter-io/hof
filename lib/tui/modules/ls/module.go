package ls

import "github.com/hofstadter-io/hof/lib/connector"

func New() connector.Connector {
	items := []interface{}{
		NewLS(),
	}

	m := connector.New("ls")
	m.Add(items)

	return m
}

