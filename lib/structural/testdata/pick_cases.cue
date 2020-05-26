package testdata

@cases(pick)
#PickCases: {
	@group(simple)
	simple: {
		t_0001: {
			args: {
				orig: #SharedExamples.A
				pick: {a: "a", N: {x: "x"}}
			}
			ex: { N: {x: "x"}, a: "a" }
		}
		t_0002: {
			args: {
				orig: #SharedExamples.A
				pick: {a: string, b: int}
			}
			ex: {a: "a"}
		}
	}
}
