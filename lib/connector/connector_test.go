package connector

import (
	"fmt"
	"testing"

	. "github.com/franela/goblin"
)

func TestNew(t *testing.T) {
	g := Goblin(t)
	g.Describe("A New Connection", func() {
		g.It("should be nil if not given a name", func() {
			g.Assert(New("")).Equal((*Base)(nil))
		})

		conn := New("name")
		g.It("should have a non-empty name", func() {
			g.Assert(conn.Name() != "").IsTrue("Name is: " + conn.Name())
		})
		g.It("should have a name equal to the New() input", func() {
			g.Assert(conn.Name()).Equal("name")
		})

		g.It("should have non-nil items", func() {
			g.Assert(conn.Items() != nil).IsTrue(fmt.Sprint("Items is nil: ", conn.Items()))
		})
		g.It("should have len(0) items", func() {
			g.Assert(len(conn.Items())).Equal(0)
		})

		g.It("should have non-nil connectors", func() {
			g.Assert(conn.Connectors() != nil).IsTrue(fmt.Sprint("Connectors is nil: ", conn.Connectors()))
		})
		g.It("should have len(0) connectors", func() {
			g.Assert(len(conn.Connectors())).Equal(0)
		})
	})
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

func TestAdd(t *testing.T) {
	g := Goblin(t)
	g.Describe("Adding to a Connection", func() {

		g.It("should be able to add a single item", func() {
			conn := New("my-connector")

			conn.Add(foo{})
			g.Assert(len(conn.Items())).Equal(1)
		})
		g.It("should be able to add multple single item", func() {
			conn := New("my-connector")

			conn.Add(foo{}, boo{})
			g.Assert(len(conn.Items())).Equal(2)
			conn.Add(foo{}, boo{})
			g.Assert(len(conn.Items())).Equal(4)
		})

		g.It("should be able to add a slice item", func() {
			conn := New("my-connector")

			conn.Add([]interface{}{foo{}, boo{}})
			g.Assert(len(conn.Items())).Equal(2)
		})
		g.It("should be able to add multiple slice item", func() {
			conn := New("my-connector")

			conn.Add([]interface{}{foo{}, boo{}})
			g.Assert(len(conn.Items())).Equal(2)

			conn.Add([]interface{}{foo{}, boo{}}, []interface{}{foo{}, boo{}, moo{}})
			g.Assert(len(conn.Items())).Equal(7)
		})

		g.It("should be able to add a mixture of items", func() {
			conn := New("my-connector")

			conn.Add(moo{}, []interface{}{foo{}, boo{}}, foo{}, []interface{}{foo{}, boo{}, moo{}})
			g.Assert(len(conn.Items())).Equal(7)
		})

	})
}

func TestGet(t *testing.T) {
	g := Goblin(t)
	g.Describe("Getting from a Connection", func() {
		conn := New("my-connector")
		conn.Add(&foo{}, &boo{}, &moo{})

		g.It("should start with multple items", func() {
			g.Assert(len(conn.Items())).Equal(3)
		})

		g.It("should get all named items", func() {
			g.Assert(len(conn.Named())).Equal(3)
		})

		g.It("should get all Named items", func() {
			g.Assert(len(conn.Get((*Named)(nil)))).Equal(3)
		})

		g.It("should get all Doer items", func() {
			g.Assert(len(conn.Get((*Doer)(nil)))).Equal(3)
		})

		g.It("should get all Talker items", func() {
			g.Assert(len(conn.Get((*Talker)(nil)))).Equal(2)
		})
	})
}
