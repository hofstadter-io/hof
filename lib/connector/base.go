package connector

import (
	"reflect"
)

// Base is the base implementation
type Base struct {
	name  string
	items []interface{}
}

// New creates a new Base (i.e. Connector) with the given name.
// Name must have a non-empty value.
// It should also be unique, but that is not enforced.
func New(name string, items ...interface{}) *Base {
	if name == "" {
		return nil
	}
	B := &Base{
		name:  name,
		items: []interface{}{},
	}

	if len(items) > 0 {
		B.Add(items)
	}

	return B
}

// Name returns the name of the Base (i.e. Connector)
func (B *Base) Name() string {
	return B.name
}

// Named returns the name of the Base (i.e. Connector)
func (B *Base) Named() []Named {
	// make a copy
	all := []Named{}
	for _, item := range B.items {
		named, ok := item.(Named)
		if ok {
			all = append(all, named)
		}
	}

	return all
}

// Items returns all Items
func (B *Base) Items() []interface{} {
	// make a copy
	all := []interface{}{}
	for _, item := range B.items {
		all = append(all, item)
	}

	return all
}

// Add adds items to a Connector. Input may be a single object, slice, or any mix.
func (B *Base) Add(in ...interface{}) {
	B.add(in)
}

func (B *Base) add(in interface{}) {
	switch it := in.(type) {
	case []interface{}:
		for _, i := range it {
			B.add(i)
		}
	default:
		B.items = append(B.items, it)

		itmzr, ok := in.(Itemizer)
		if ok {
			for _, i := range itmzr.Items() {
				B.items = append(B.items, i)
			}
		}
	}
}

// Del deletes an item or list of items from the connector.
func (B *Base) Del(out interface{}) {
	switch ot := out.(type) {
	case []interface{}:
		for _, o := range ot {
			B.Del(o)
		}
		return
	}

	// TODO otherwise look for objects
}

// Clear clears all items from this Connector.
// Will not effect connectors which this had been added to.
func (B *Base) Clear() {
	B.items = []interface{}{}
}

/*Get will return items which implement a given type or interface.

Once a number of items have been added,
we will want to retrieve some subset of those.
The Get() Function will return all items
that match the given type
*/
func (B *Base) Get(in interface{}) []interface{} {
	typ := reflect.TypeOf(in).Elem()

	// make a copy
	all := []interface{}{}
	for _, item := range B.items {
		it := reflect.TypeOf(item)
		if it.Implements(typ) {
			all = append(all, item)
		}
	}

	return all
}

/*Connect is the main usage function.

Connect() recursively passes a Connector object to all items
so that they may consume or use any items in the Connector.
Typically, we build up a root Connector and then
call Connect with itself as the argument.
*/
func (B *Base) Connect(c Connector) {
	for _, item := range B.Items() {
		switch typ := item.(type) {
		case Connectable:
			typ.Connect(c)
		}
	}
}
