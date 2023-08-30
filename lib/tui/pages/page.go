package pages

import (
	"github.com/rivo/tview"
)

type Page struct {
	Name string

	*tview.Flex
}
