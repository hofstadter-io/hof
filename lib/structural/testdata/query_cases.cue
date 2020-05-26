package testdata

@cases(query)
#QueryCases: {
	@group(simple)
	simple: {
		t_0001: {
			args: {
				orig: #SharedExamples.A
				query: {a: string, N: {x: "x"}}
			}
			ex: {a: "a", N: {x: "x"}}
		}
		t_0002: {
			args: {
				orig: #SharedExamples.A
				query: {a: string, b: int}
			}
			ex: {a: "a"}
		}
	}
}
