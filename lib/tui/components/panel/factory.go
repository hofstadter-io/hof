package panel

import (
	"fmt"
)

// I hope this doesn't trigger you :]
// It is a default helper for creating interesting CUE widgets and quick actions for them
// note, you should always add a "default" key, as this is used in panel when CRUD'n items/nested panels
type Factory struct {

	makers map[string]ItemCreator

}

func NewFactory() *Factory{
	return &Factory{
		makers: make(map[string]ItemCreator),
	}
}

func (F *Factory) Register(itemKey string, creator ItemCreator) {
	F.makers[itemKey] = creator
}

func (F *Factory) Creator(context ItemContext, parent *Panel) (PanelItem, error) {
	// tui.Log("debug", context)

	// cleanup args, loading json can create []any when restoring a panel or item
	args := []string{}
	if _args, ok := context["args"]; ok {
		// because in-mem vs decode-yaml...
		switch _args := _args.(type) {
		case []string:
			args = _args
		case []any:
			for _, a := range _args {
				args = append(args, a.(string))
			}
		}
	}
	context["args"] = args

	//
	// this should all go in a cue item and a cue creator function
	//   eval would then not know about these and
	//   use the cue creator as the arg to dash
	//

	// get the item type from context
	item := ""
	if _item, ok := context["item"]; ok {
		item = _item.(string)
	}


	maker, ok := F.makers[item]
	if !ok {
		return nil, fmt.Errorf("unknown creator: %q", item)
	}

	return maker(context, parent)
}

