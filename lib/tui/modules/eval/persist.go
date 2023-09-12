package eval

import (
	"fmt"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/playground"
	"github.com/hofstadter-io/hof/lib/tui/components/panel"
	"github.com/hofstadter-io/hof/lib/tui/components/widget"
)

const evalSaveDirSubdir = "tui/saves/eval"

func evalSavePath(filename string) string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir,"hof",evalSaveDirSubdir, filename)
}

func (M *Eval) Save(destination string, preview bool) error {

	//
	// encode and marshal the dashboard
	//
	m, err := M.EncodeMap()
	if err != nil {
		return err
	}

	ctx := cuecontext.New()
	v := ctx.CompileString("{}")
	v = v.FillPath(cue.ParsePath(""), m)

	//b, err := yaml.Marshal(m)
	//if err != nil {
	//  return err
	//}


	if preview {

		cfg := &helpers.SourceConfig{Value: v, Source: helpers.EvalNone}
		t := browser.New(cfg, "cue")
		t.Rebuild()

		//t := widget.NewTextView()
		//t.SetDynamicColors(false)
		//fmt.Fprint(t, string(b))
		I := panel.NewBaseItem(M.Panel)
		I.SetWidget(t)
		M.AddItem(I, 0, 1, true)

	} else {
		// save location
		savename := evalSavePath(destination)

		opts := []cue.Option{
			cue.Final(),
		}
		syn := v.Syntax(opts...)

		b, err := format.Node(syn)
		if err != nil {
			return err
		}	

		// ensure the dir exists
		dir := filepath.Dir(savename)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		// write our dashboard out
		err = os.WriteFile(savename, b, 0644)
		if err != nil {
			return err
		}

		//
		// alert the user in several ways
		//

		info := fmt.Sprintf("%s saved to ... %s", M.Name(), savename)
		tui.Tell("info", info)
		tui.Log("info", info)
	}

	return nil
}

func (M *Eval) LoadEval(filename string) (*Eval, error) {
	tui.Log("debug", fmt.Sprintf("Eval.LoadEval.0: %v", filename))
	savename := evalSavePath(filename)

	b, err := os.ReadFile(savename)
	tui.Log("debug", fmt.Sprintf("Eval.LoadEval.1: %v %v %v", savename, len(b), err))
	if err != nil {
		return nil, err
	}

	ctx := cuecontext.New()
	val := ctx.CompileBytes(b, cue.Filename(filename))

	data := make(map[string]any)
	err = val.Decode(&data)
	if err != nil {
		return nil, err
	}

	// todo, decode the actual thing
	e, err := EvalDecodeMap(data)
	if err != nil {
		return nil, err
	}

	M.Panel = e.Panel
	M.showPanel = e.showPanel
	M.showOther = e.showOther
	M.Panel.SetTitle(M.Panel.TitleString())

	tui.SetFocus(M.Panel)

	// M.Mount(make(map[string]any))


	// extra to display the save info
	//t := NewTextView()
	//t.SetDynamicColors(false)
	//fmt.Fprint(t, string(b))
	//I := NewItem(nil, M.Panel)
	//I.SetWidget(t)
	//M.AddItem(I, 0, 1, true)

	return nil, nil
}

func (M *Eval) ShowEval(filename string) (*Eval, error) {
	savename := evalSavePath(filename)

	b, err := os.ReadFile(savename)
	if err != nil {
		return nil, err
	}

	// extra to display the save info
	t := widget.NewTextView()
	t.SetDynamicColors(false)
	fmt.Fprint(t, string(b))
	I := panel.NewBaseItem(M.Panel)
	I.SetWidget(t)
	M.AddItem(I, 0, 1, true)

	return nil, nil
}

func (M *Eval) ListEval() (error) {
	dir := evalSavePath("")

	infos, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// start our new text view
	t := widget.NewTextView()
	t.SetDynamicColors(false)

	// write listing
	for _, info := range infos {
		fmt.Fprintln(t.TextView, info.Name())
	}

	// display the file list to the user
	I := panel.NewBaseItem(M.Panel)
	I.SetWidget(t)
	M.AddItem(I, 0, 1, true)

	return nil
}

func (M *Eval) EncodeMap() (map[string]any, error) {
	var err error
	m := make(map[string]any)

	// metadata
	m["name"] = M.Name()
	m["type"] = "eval"

	// visual settings
	m["direction"] = M.GetDirection()
	m["showPanel"] = M.showPanel
	m["showOther"] = M.showOther
	
	// panel
	m["panel"], err = M.Panel.Encode()
	if err != nil {
		return m, err
	}

	return m, nil
}

func EvalDecodeMap(input map[string]any) (*Eval, error) {
	// tui.Log("extra", fmt.Sprintf("Eval.Decode: %# v", input))
	var err error

	M := &Eval{
		showPanel: input["showPanel"].(bool),
		showOther: input["showOther"].(bool),
	}
	pmap, ok := input["panel"]
	if !ok {
		return nil, fmt.Errorf("panel missing from eval decode input: %#v", input)
	} else {
		M.Panel = panel.New(nil, M.creator)
	}



	// decode the main panel, everything else should happen through recursion and widget registry
	M.Panel, err = panel.PanelDecodeMap(pmap.(map[string]any), nil, M.creator)
	if err != nil {
		return M, err
	}

	M.Panel.SetTitle(M.Panel.TitleString())

	// do layout setup here, once some children have been instantiated
	M.SetName(input["name"].(string))
	M.SetDirection(input["direction"].(int))
	M.SetBorderColor(tcell.Color42).SetBorder(true)
	M.SetBorder(true)
	M.setupEventHandlers()

	// M.Refresh(make(map[string]any))

	// tui.Draw()

	err = M.reconnectItems()

	return M, err
}

func (M *Eval) reconnectItems() error {

	reconnect := func (p panel.PanelItem) {
		switch t := p.Widget().(type) {
		case *playground.Playground:
			dstPlay := t
			sc := dstPlay.GetScopeConfig()
			if sc.Source == helpers.EvalConn {
				args := sc.Args

				path := args[0]
				srcItem, err := M.getItemByPath(path)
				if err != nil {
					tui.Log("error", err)
					return
				}

				var expr string
				if len(args) == 2 {
					expr = args[1]
				}

				switch t := srcItem.Widget().(type) {
				case *playground.Playground:
					if expr == "" {
						dstPlay.SetConnection(args, t.GetConnValue)
					} else {
						dstPlay.SetConnection(args, t.GetConnValueExpr(expr))
					}
				case *browser.Browser:
					if expr == "" {
						dstPlay.SetConnection(args, t.GetConnValue)
					} else {
						dstPlay.SetConnection(args, t.GetConnValueExpr(expr))
					}
				}
			}
		}
	}

	M.Panel.RangeItems(reconnect)

	return nil
}
