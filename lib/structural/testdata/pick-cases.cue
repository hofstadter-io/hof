pick_cases: [
	{
		value: {
			a: "a"
			b: "b"
			d: "d"
			e: {
				a: "a"
				b: "b1"
				d: "cd"
			}
		}
		pick: {
			b: string
			d: int
			e: {
				a: _
				b: =~"^b"
				d: =~"^d"
			}
    }
		expect: {
			b: "b"
			e:  {
				a: "a"
				b: "b1"
			}
		}
	},
	{
		value: {
			c: {
				a: "a"
				b: "b"
				c: {
					a: "a"
					b: "b"
				}
			}
		}
		pick: {
			c: {
				a: _
			}
    }
		expect: {
			c:  {
				a: "a"
			}
		}
	},
	// {
	// 	value: { string }
	// 	pick: { "foo" }
	// 	expect: {
	// 		"foo"
	// 	}
	// },
]
