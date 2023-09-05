package eval

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/hofstadter-io/hof/lib/tui"
)


const evalSaveDirSubdir = "tui/saves/eval"

func evalSavePath(filename string) string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir,"hof",evalSaveDirSubdir, filename)
}

func (M *Eval) Save(filename string) error {

	//
	// encode and marshal the dashboard
	//
	m, err := M.EncodeMap()
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	// save location
	savename := evalSavePath(filename)

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

	// extra to display the save info
	//t := NewTextView()
	//t.SetDynamicColors(false)
	//fmt.Fprint(t, string(b))
	//I := NewItem(nil, M.Panel)
	//I.SetWidget(t)

	//M.Panel.AddItem(I, 0, 1, true)

	return nil
}

func (M *Eval) LoadEval(filename string) (*Eval, error) {
	savename := evalSavePath(filename)

	b, err := os.ReadFile(savename)
	tui.Log("debug", fmt.Sprintf("Eval.LoadEval.1: %v %v %v", savename, len(b), err))
	if err != nil {
		return nil, err
	}

	data := make(map[string]any)
	err = yaml.Unmarshal(b, &data)
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
	M._cnt = e._cnt
	M._name = e._name


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
	tui.Log("debug", fmt.Sprintf("Eval.LoadEval.1: %v %v %v", savename, len(b), err))
	if err != nil {
		return nil, err
	}

	// extra to display the save info
	t := NewTextView()
	t.SetDynamicColors(false)
	fmt.Fprint(t, string(b))
	I := NewItem(nil, M.Panel)
	I.SetWidget(t)
	M.AddItem(I, 0, 1, true)

	return nil, nil
}

func (M *Eval) ListEval() (error) {
	dir := evalSavePath("")
	tui.Log("debug", fmt.Sprintf("Eval.ListEval.1: %v", dir))

	infos, err := os.ReadDir(dir)
	tui.Log("debug", fmt.Sprintf("Eval.ListEval.2: %v %v", len(infos), err))
	if err != nil {
		return err
	}

	// start our new text view
	t := NewTextView()
	t.SetDynamicColors(false)

	// write listing
	for _, info := range infos {
		fmt.Fprintln(t.TextView, info.Name())
	}

	// display the file list to the user
	I := NewItem(nil, M.Panel)
	I.SetWidget(t)
	M.AddItem(I, 0, 1, true)

	return nil
}

func (M *Eval) EncodeMap() (map[string]any, error) {
	var err error
	m := make(map[string]any)

	// metadata
	m["id"] = M._cnt
	m["name"] = M._name
	m["type"] = "eval"

	// visual settings
	m["direction"] = M.GetDirection()
	m["showPanel"] = M.showPanel
	m["showOther"] = M.showOther
	
	// panel
	m["panel"], err = M.Panel.EncodeMap()
	if err != nil {
		return m, err
	}

	return m, nil
}

func EvalDecodeMap(data map[string]any) (*Eval, error) {
	var err error
	M := &Eval{
		showPanel: data["showPanel"].(bool),
		showOther: data["showOther"].(bool),
		_cnt: data["id"].(int),
		_name: data["name"].(string),
	}

	if pmap, ok := data["panel"]; ok {
		M.Panel, err = PanelDecodeMap(pmap.(map[string]any), nil, nil)
		if err != nil {
			return M, err
		}
	} else {
		M.Panel = NewPanel(nil, nil)
	}

	// do layout setup here, once some children have been instantiated
	M.Flex.SetDirection(data["direction"].(int))
	M.Flex.SetBorder(true).SetTitle(M.Name())

	return M, nil
}
