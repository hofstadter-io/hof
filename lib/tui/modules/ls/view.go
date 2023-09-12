package ls

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/common"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// Both a Module and a Layout and a Switcher.SubLayout
type LS struct {
	*tview.Flex

	Tree *common.FileBrowser
	View *common.TextEditor
}

func NewLS() *LS {
	page := &LS{
		Flex:	tview.NewFlex(),
	}

	return page
}

func (P *LS) Id() string {
	return "ls"
}

func (P *LS) Mount(context map[string]any) error {
	P.Flex.Mount(context)

	return P.Refresh(context)
}

func (P *LS) Refresh(context map[string]any) error {
	dir := ""
	_args, _ := context["args"]
	args, _ := _args.([]string)
	if len(args) > 0 && args[0] == "ls" {
		args = args[1:]
	}

	if len(args) > 0 {
		dir = args[0]
	} else {
		dir, _ = os.Getwd()
	}
	dir, _ = filepath.Abs(dir)
	go tui.SendCustomEvent("/console/info", "LS dir: " + dir)
	
	P.View = common.NewTextEditor(nil)

	P.Tree = common.NewFileBrowser(dir,
		func(path string) {
			P.View.OpenFile(path)
		},
		func(path string) {
			wd, _ := os.Getwd()
			rpath := strings.TrimPrefix(path, wd + "/")
			context := make(map[string]any)
			context["args"] = []string{rpath}
			context["path"] = "/eval"
			context["mode"] = "code"
			go tui.SendCustomEvent("/router/dispatch", context)
		},
	)

	P.Flex.Clear()
	P.AddItem(P.Tree, 42, 1, false)
	P.AddItem(P.View, 0, 1, false)

	P.Tree.Mount(context)
	P.View.Mount(context)

	P.Tree.Focus(nil)

	return nil
}

func (P *LS) Name() string {
	return "ls"
}

func (P *LS) HotKey() string {
	return ""
}

func (P *LS) CommandName() string {
	return "ls"
}

func (P *LS) CommandUsage() string {
	return "ls [path]"
}

func (P *LS) CommandHelp() string {
	return "show the path in the file browser"
}

func (P *LS) CommandCallback(context map[string]any) {
	if P.IsMounted() {
		// just refresh with new args
		P.Refresh(context)
	} else {
		// need to navigate, mount will do the rest
		context["path"] = "/ls"
		go tui.SendCustomEvent("/router/dispatch", context)
	}
}
