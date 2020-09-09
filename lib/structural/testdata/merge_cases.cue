package testdata

@cases(merge)
#MergeCases: {
	@group(simple)
	edge: {
		t_0001: {
			args: {
				orig: {a: "a"}
				merge: {}
			}
			ex: {a: "a"}
		}
		t_0002: {
			args: {
				orig: {}
				merge: {a: "a"}
			}
			ex: {a: "a"}
		}
	}
	//@group(simple)
	//simple: {
		//t_0002: {
			//args: {
				//orig: #SharedExamples.A
				//merge: {a: "a", N: {x: "x"}}
			//}
			//ex: {a: "a", b: "b", N: {x: "x", y: "y"}}
		//}
	//}
}
