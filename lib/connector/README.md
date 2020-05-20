# connector-go

[![Build Status](https://travis-ci.org/hofstadter-io/connector-go.svg?branch=master)](https://travis-ci.org/hofstadter-io/connector-go)
[![Doc Status](https://godoc.org/github.com/hofstadter-io/connector-go?status.png)](https://godoc.org/github.com/hofstadter-io/connector-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/hofstadter-io/connector-go)](https://goreportcard.com/report/github.com/hofstadter-io/connector-go)


A Golang implementation of the __Connector__ concept.

Examples:

- https://github.com/verdverm/vermui-starterkit

Basic Code:

```go
import "github.com/hofstadter-io/connector-go"

func main () {
    conn := connector.New("my-connector")
    f, b, m := foo{do:"goo"},boo{do:"be friendly"},moo{do:"farm to table"}
    conn.Add(f, []interface{}{b,m})

    for _, named := range conn.Named() {
        named.Name()
    }

    for _, item := range conn.Items() {
        doer, ok := item.(Doer)
        if ok {
            doer.Do()
        }
    }

    for _, item := range conn.Get((*Talker)(nil)) {
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
```
