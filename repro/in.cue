package foo

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

@flow()

tasks: {

	p1: {val: x, pick: y} @task(st.Pick) @print(out)

	m1: {val: p1.out, mask: {c: int}} @task(st.Mask) @print(out)
	m2: {val: p1.out, mask: {a: _}} @task(st.Mask) @print(out)

	u1: {val: m1.out, upsert: m2.out} @task(st.Upsert) @print(out)
}
