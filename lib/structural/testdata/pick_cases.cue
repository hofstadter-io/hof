package testdata

@cases(pick)
#PickCases: {
	@group(original)
	original: {
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

	@group(simple)
	simple: {
		t_0001: {
			args: {
				orig: #SharedExamples.a1
				pick: {a: "a"}
			}
			ex: { a: "a" }
		}
		t_0002: {
			args: {
				orig: #SharedExamples.a1
				pick: {a: string}
			}
			ex: {a: "a"}
		}
		t_0003: {
			args: {
				orig: #SharedExamples.a1
				pick: {a: =~ "[a-z]+"}
			}
			ex: {a: "a"}
		}
		t_0004: {
			args: {
				orig: #SharedExamples.a1
				pick: {a: _}
			}
			ex: {a: "a"}
		}
		//neg_diffval: {
			//args: {
				//orig: #SharedExamples.a1
				//pick: {a: "b"}
			//}
			//ex: {}
		//}
		neg_difftype_1: {
			args: {
				orig: #SharedExamples.a1
				pick: {a: int}
			}
			ex: {}
		}
		neg_difftype_1: {
			args: {
				orig: #SharedExamples.a1
				pick: {a: number}
			}
			ex: {}
		}
		neg_diffregexp: {
			args: {
				orig: #SharedExamples.a1
				pick: {a: =~ "[A-Z]+"}
			}
			ex: {}
		}
	}
}
