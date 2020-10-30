package testdata

#PickCases: {

	basiclit: {
		"str check any": {
			args: {
				orig: a: "str"
				pick: a: _
			}
			ex: {
				a: "str"
			}
		}
		"str check same": {
			args: {
				orig: a: "str"
				pick: a: "str"
			}
			ex: {
				a: "str"
			}
		}
		"str check type": {
			args: {
				orig: a: "str"
				pick: a: string
			}
			ex: {
				a: "str"
			}
		}
		"str check regex": {
			args: {
				orig: a: "str"
				pick: a: =~ "[a-z]+"
			}
			ex: {
				a: "str"
			}
		}
		"str check no unify": {
			args: {
				orig: a: "str"
				pick: a: int
			}
			ex: {}
		}
	}

	structs: {
		"check any": {
			args: {
				orig: a: { b: "b" }
				pick: a: _
			}
			ex: {
				a: { b: "b" }
			}
		}
		"struct check closed-open struct": {
			args: {
				orig: a: { b: "b" }
				pick: a: {...}
			}
			ex: {
				a: { b: "b" }
			}
		}
		"struct check open-open struct": {
			args: {
				orig: a: { b: "b", ... }
				pick: a: {...}
			}
			ex: {
				a: { b: "b" }
			}
		}
		"struct check closed-closed struct": {
			args: {
				orig: a: { b: "b" }
				pick: a: close({})
			}
			ex: {}
		}
		"struct check open-closed struct": {
			args: {
				orig: a: { b: "b", ... }
				pick: a: close({})
			}
			ex: {}
		}
		"check struct with str, should be empty": {
			args: {
				orig: a: { b: "b" }
				pick: a: "a"
			}
			ex: {}
		}
	}

	//simple: {
		//t_0001: {
			//args: {
				//orig: #SharedExamples.a1
				//pick: {a: "a"}
			//}
			//ex: { a: "a" }
		//}
		//t_0002: {
			//args: {
				//orig: #SharedExamples.a1
				//pick: {a: string}
			//}
			//ex: {a: "a"}
		//}
		//t_0003: {
			//args: {
				//orig: #SharedExamples.a1
				//pick: {a: =~ "[a-z]+"}
			//}
			//ex: {a: "a"}
		//}
		//t_0004: {
			//args: {
				//orig: #SharedExamples.a1
				//pick: {a: _}
			//}
			//ex: {a: "a"}
		//}
		//neg_diffval: {
			//args: {
				//orig: #SharedExamples.a1
				//pick: {a: "b"}
			//}
			//ex: {}
		//}
		//neg_difftype_1: {
			//args: {
				//orig: #SharedExamples.a1
				//pick: {a: int}
			//}
			//ex: {}
		//}
		//neg_difftype_1: {
			//args: {
				//orig: #SharedExamples.a1
				//pick: {a: number}
			//}
			//ex: {}
		//}
		//neg_diffregexp: {
			//args: {
				//orig: #SharedExamples.a1
				//pick: {a: =~ "[A-Z]+"}
			//}
			//ex: {}
		//}
	//}

	//nested: {
		//t_0001: {
			//args: {
				//orig: #SharedExamples.A
				//pick: {a: "a", N: {x: "x"}}
			//}
			//ex: { a: "a", N: {x: "x"} }
		//}
		//t_0002: {
			//args: {
				//orig: #SharedExamples.A
				//pick: {a: string, b: int}
			//}
			//ex: {a: "a"}
		//}
	//}

}
