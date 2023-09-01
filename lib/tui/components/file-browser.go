package components

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui/app"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type FileBrowser struct {
	Dir string

	OnOpen func(string)

	*tview.TreeView

	Root *tview.TreeNode
	Node *tview.TreeNode

	App *app.App
}


func NewFileBrowser(app *app.App, dir string, onopen func(path string)) *FileBrowser {
	fb := &FileBrowser {
		App: app,
		Dir: dir,
		OnOpen: onopen,
	}

	// file browser
	fb.Root = tview.NewTreeNode(dir)
	fb.Root.SetColor(tcell.ColorRed)
	fb.AddAt(fb.Root, dir)

	// tree view
	fb.TreeView = tview.NewTreeView()
	fb.
		SetRoot(fb.Root).
		SetCurrentNode(fb.Root)
	fb.SetBorder(true)

	// set our selected handler
	fb.SetSelectedFunc(fb.OnSelect)


	return fb
}

func (FB *FileBrowser) OnSelect(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}

	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		path := reference.(string)
		info, _ := os.Lstat(path)
		if info.IsDir() {
			FB.AddAt(node, path)
		} else {
			FB.OnOpen(path)
		}
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}


func (FB *FileBrowser) AddAt(target *tview.TreeNode, path string) {
	// get files at path
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	// sort dirs first, then by name
	sort.Slice(files, func(x, y int) bool {
		X, Y := files[x], files[y]
		// deal with file vs dir
		if X.IsDir() && !Y.IsDir() {
			return true
		} else if !X.IsDir() && Y.IsDir() {
			return false
		} else {
			return X.Name() < Y.Name()
		}
	})

	// build tree nodes
	for _, file := range files {
		node := tview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(path, file.Name()))
			// SetSelectable(file.IsDir())
		if file.IsDir() {
			node.SetColor(tcell.ColorGreen)
		}
		target.AddChild(node)
	}
}
