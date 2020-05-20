/*
Package connector ...

Implements the Connector concept in Golang.


	import "github.com/hofstadter-io/connector-go"

	func main () {
		conn := connector.New("my-connector")
		f, b, m := foo{do:"goo"},boo{do:"be friendly"},moo{do:"farm to table"}
		conn.Add(f, []interface{}{b,m})

		for _, named := range conn.Named() {
			named.Name()
		}

		for _, item := range conn.Items() {
			doer := item.(Doer)
			doer.Do()
		}

		typ := reflect.TypeOf((*Talker)(nil)).Elem()
		for _, item := range conn.Get(typ) {
			talker := item.(Talker)
			talker.Say()
		}
	}


	type Doer interface {
		Do() string
	}

	type Talker interface {
		Say() string
	}

	type foo struct {
		do string
	}

	func (f *foo) Do() string {
		return f.do
	}
	func (f *foo) Name() string {
		return "foo"
	}

	type boo struct {
		do string
	}

	func (b *boo) Do() string {
		return b.do
	}
	func (b *boo) Name() string {
		return "Casper"
	}
	func (b *boo) Say() string {
		return "Boooooo"
	}

	type moo struct {
		do string
	}

	func (m *moo) Do() string {
		return m.do
	}
	func (m *moo) Name() string {
		return "Cow"
	}
	func (m *moo) Say() string {
		return "MoooOOO"
	}
*/
package connector
