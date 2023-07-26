package main

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"

	st "github.com/hofstadter-io/hof/lib/structural"
)


const code =`
x: {
	"a": {
		"b": "B"
	}
	"b": 1
	"c": 2
	"d": "D"
}

y: {
	a: {
		b: string
	}
	c: int
	d: "D"
}

// x pick y
p1: {
	a: {
		b: "B"
	}
	c: 2
	d: "D"
}

// p1 mask m1p
m1: {
	a: {
		b: "B"
	}
	d: "D"
}
// p1 mask m2p
m2: {
	c: 2
	d: "D"
}

m1p: { c: int }
m2p: { a: _ }
`

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	c := cuecontext.New()

	// read and compile value
	val := c.CompileString(code)
	fmt.Println("val:", val)

	x := val.LookupPath(cue.ParsePath("x"))
	y := val.LookupPath(cue.ParsePath("y"))

	// p1 := val.LookupPath(cue.ParsePath("p1"))
	p1, err := st.PickValue(y, x, nil)
	checkErr(err)
	fmt.Println("p1:", p1)

	// Swapping here, it works
	//m1 := val.LookupPath(cue.ParsePath("m1"))
	m1p := val.LookupPath(cue.ParsePath("m1p"))
	m1, err := st.MaskValue(m1p, p1, nil)
	checkErr(err)
	fmt.Println("m1:", m1)

	m2p := val.LookupPath(cue.ParsePath("m2p"))
	m2, err := st.MaskValue(m2p, p1, nil)
	checkErr(err)
	fmt.Println("m2:", m2)

	// swapping m1, m2 will produce a `c: _` which does not exist anywhere
	u1, err := st.UpsertValue(m2, m1, nil)
	checkErr(err)
	fmt.Println("u1:", u1)
}
