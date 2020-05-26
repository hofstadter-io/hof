package testdata

@cases(mask)
#MaskCases: {
	@group(simple)
	simple: {
		t_0001: {
			args: {
				orig: #SharedExamples.A
				mask: {a: string, N: {x: "x"}}
			}
			ex: {b: "b", N: {y: "y"}}
		}
		t_0002: {
			args: {
				orig: #SharedExamples.A
				mask: {a: string, b: int}
			}
			ex: {b: "b", N: {x: "x", y: "y"}}
		}
	}
}
