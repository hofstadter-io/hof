package common

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/dotpath"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/widget"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type MapBrowser struct {
	Name string
	Data map[string]any

	OnOpen  func(string)
	OnClick func(string)

	LeafClick func(string)

	*tview.TreeView

	Root *tview.TreeNode
	Node *tview.TreeNode
}


func NewMapBrowser(name string, data map[string]any, onopen, onclick func(path string)) *MapBrowser {

	C := &MapBrowser {
		Name: name,
		Data: data,
		OnOpen: onopen,
		OnClick: onclick,
	}

	// file browser
	C.Root = tview.NewTreeNode(C.Name)
	C.Root.SetColor(tcell.ColorAqua)
	C.Root.SetSelectable(true)
	C.Root.SetDoubleClickFunc(func(){

		// this does not seem to work...
		tui.Log("crit", "root double click")
		if C.Root.IsExpanded() {
			C.Root.CollapseAll()
		} else {
			C.Root.ExpandAll()
		}
	})
	C.AddAt(C.Root, C.Name)

	// tree view
	C.TreeView = tview.NewTreeView()
	C.
		SetRoot(C.Root).
		SetCurrentNode(C.Root)
	C.SetBorder(true)

	// set our selected handler
	C.SetSelectedFunc(C.OnSelect)
	C.SetDoubleClickFunc(C.OnLeafClick)

	return C
}

func (C *MapBrowser) OnLeafClick(node *tview.TreeNode) {
	tui.Log("debug", "map-tree ... click!")
	if C.LeafClick == nil {
		return
	}

	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}

	path := reference.(string)

	tui.Log("debug", "map-tree ... on " + path)
	C.LeafClick(path)
}

func (C *MapBrowser) OnSelect(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}

	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		path := reference.(string)
		path = strings.ReplaceAll(path, "/", ".")
		d, _ := dotpath.Get(path, C.Data, false)
		if d != nil {
			C.AddAt(node, path)
		}
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}


func (C *MapBrowser) AddAt(target *tview.TreeNode, path string) {
	// get data at path
	d, _ := dotpath.Get(path, C.Data, false)

	switch m := d.(type) {
	case map[string]any:
		// build tree nodes
		for k, v := range m {
			node := tview.NewTreeNode(k).SetReference(filepath.Join(path, k))
			switch v.(type) {
			case map[string]any, []any, []string:
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	case []any:
		for i, v := range m {
			s := fmt.Sprintf("[%d]", i)
			node := tview.NewTreeNode(s).SetReference(path + s)
			switch v.(type) {
			case map[string]any, []any, []string:
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	case []string:
		for i, v := range m {
			s := fmt.Sprintf("[%d]", i)
			node := tview.NewTreeNode(v).SetReference(path + s)
			target.AddChild(node)
		}
	}
}

func (C *MapBrowser) TypeName() (string) {
	return "common/map-browser"
}

func (C *MapBrowser) Encode() (map[string]any, error) {
	return nil, nil
}

func (C *MapBrowser) Decode(map[string]any) (widget.Widget, error) {
	return nil, nil
}
