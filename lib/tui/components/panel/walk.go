package panel


func (P* Panel) RangeItems(fn func (PanelItem)) {

	for _, item := range P.GetItems() {
		fi := item.Item // flexItem.Item
		switch t := fi.(type) {
		case *Panel:
			t.RangeItems(fn)
		case PanelItem:
			fn(t)
		}
	}
}
