package eval

import (
	"encoding/json"
	"fmt"

	"github.com/hofstadter-io/hof/lib/tui"
)



func (M *Eval) Save(filename string) error {

	m, err := M.EncodeMap()
	if err != nil {
		tui.Tell("error", err)
		tui.Log("error", err)
		return err
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		tui.Tell("error", err)
		tui.Log("error", err)
		return err
	}

	t := NewTextView()
	t.SetDynamicColors(false)
	fmt.Fprint(t, string(b))
	I := NewItem(nil, M.Panel)
	I.SetWidget(t)

	M.Panel.AddItem(I, 0, 1, true)

	info := fmt.Sprintf("%s saved to ... %s", M.Name(), filename)
	tui.Tell("info", info)
	tui.Log("info", info)

	return nil
}

func LoadEval(filename string) (*Eval, error) {

	return nil, nil
}


func (M *Eval) EncodeMap() (map[string]any, error) {
	var err error
	m := make(map[string]any)

	// metadata
	m["id"] = M._cnt
	m["name"] = M._name
	m["type"] = "eval"

	// visual settings
	m["showPanel"] = M.showPanel
	m["showOther"] = M.showOther
	
	// panel
	m["panel"], err = M.Panel.EncodeMap()
	if err != nil {
		return m, err
	}

	return m, nil
}

