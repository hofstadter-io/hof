package ls

import (
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components"
	"github.com/hofstadter-io/hof/lib/tui/hoc/router"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// Both a Module and a Layout and a Switcher.SubLayout
type LS struct {
	*tview.Flex

	Tree *components.FileBrowser
	View *components.TextEditor
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

func (P *LS) Routes() []router.RoutePair {
	return []router.RoutePair{
		{ Path: "/ls", Data: P },
	}
}

func (P *LS) Mount(context map[string]any) error {
	P.Flex.Mount(context)

	return P.Refresh(context)
}

func (P *LS) Refresh(context map[string]any) error {
	
	P.View = components.NewTextEditor(nil)

	P.Tree = components.NewFileBrowser("",
		func(path string) {
			P.View.OpenFile(path)
		},
		func(path string) {
			context := make(map[string]any)
			context["args"] = []string{path}
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

func (P *LS) CommandCallback(args []string, context map[string]any) {
	if context == nil {
		context = make(map[string]any)
	}
	context["args"] = args

	if P.IsMounted() {
		// just refresh with new args
		P.Refresh(context)
	} else {
		// need to navigate, mount will do the rest
		context["path"] = "/ls"
		go tui.SendCustomEvent("/router/dispatch", context)
	}
}
