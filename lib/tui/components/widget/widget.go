package widget

import (
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// base and wrapped tview widgets, temporarily here

// Widget is designed to fit in containers and be serializable
type Widget interface {
	tview.Primitive

	TypeName() string

	Encode() (map[string]any, error)
	// Decode(map[string]any) (error)

	// UpdateValue()
}
