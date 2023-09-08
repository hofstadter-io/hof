package panel

//import (
//  "fmt"

//  "github.com/hofstadter-io/hof/lib/tui/tview"
//)

func (P *Panel) Encode() (map[string]any, error) {
	m := make(map[string]any)

	m["id"] = P._cnt
	m["name"] = P._name
	m["type"] = "panel"
	m["direction"] = P.Flex.GetDirection()

	items := []map[string]any{}

	for _, item := range P.GetItems() {
		var (
			i map[string]any
			err error
		)

		switch item := item.Item.(type) {
		case *Panel:
			// recursion, within the current panel stack
			i, err = item.Encode()
			if err != nil {
				return m, err
			}
		case PanelItem:
			// leaf, but may also recurse if a widget uses panels
			i, err = item.Encode()
			if err != nil {
				return m, err
			}

		default:
			panic("unhandled item type in panel")	
		}

		// add the item to output
		items = append(items, i)
	}

	m["items"] = items

	return m, nil
}

//func PanelDecodeMap(data map[string]any, parent *Panel, creator ItemCreator) (*Panel, error) {
//  P := &Panel{
//    Flex: tview.NewFlex(),
//    _creator: creator,
//    _parent: parent,
//    _cnt: data["id"].(int),
//    _name: data["name"].(string),
//  }

//  if items, ok := data["items"]; ok {
//    for _, idata := range items.([]any) {
//      imap := idata.(map[string]any)
//      I, err := ItemDecodeMap(imap, P)
//      if err != nil {
//        return P, err
//      }
//      P.AddItem(I, 0, 1, true)
//    }
//  } else {
//    txt := NewTextView()
//    fmt.Fprint(txt, fmt.Sprintf("unhandled panel decode: \n%# v\n\n", data))
//    fmt.Fprint(txt, EvalHelpText)
//    I := NewItem(nil, parent)
//    I.SetWidget(txt)
//    P.AddItem(I, 0, 1, true)

//  }

//  // do layout setup here, once some children have been instantiated
//  P.Flex.SetDirection(data["direction"].(int))
//  P.Flex.SetBorder(true).SetTitle(P.Name())

//  return P, nil
//}


