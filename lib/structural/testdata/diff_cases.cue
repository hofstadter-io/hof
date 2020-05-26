package testdata

@cases(diff)
#DiffCases: {
	@group(simple)
	simple: {
		t_0001: {
			args: {
				orig: #SharedExamples.A
				next: {a: string, N: {x: "x"}}
			}
			ex: {
				removed: {
					b: "b"
				}
				inplace: {
					N: {
						removed: {
							y: "y"
						}
					}
				}
			}
		}
		t_0002: {
			args: {
				orig: #SharedExamples.A
				next: {a: string, b: int}
			}
			ex: {
				removed: {
					N: {
						x: "x"
						y: "y"
					}
				}
				changed: {
					b: {
						from: "b"
						to:   int
					}
				}
			}
		}
	}
}
