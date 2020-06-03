package test

import "list"

#Suites: [S=string]: [T=string]: #Tester & {
	name: "\(S)/\(T)"
	...
}

#Tester: {
	name: string
	type: string

	pass: bool | *false

	dir: string
	cmd: string

	sonar?: string
	bench?: string

	...
}

Suites: #Suites & {

	st: {
		cli: {
			pass: true
			type: "testsuite"
			dir: "lib/structural"
			cmd: "go test ./"
		}

		unit: {
			type: "gotest"
			dir: "lib/structural"
			cmd: "go test ./"
		}
	}


	mod: {
		cli: {
			type: "testsuite"
			dir: "lib/mod"
			cmd: "go test ./"
		}

		unit: {
			pass: true
			type: "gotest"
			dir: "lib/mod"
			cmd: "go test ./"
		}
	}

}

#Flags: {
	suite: string | *"all" @tag(suite,short=st|mod)
	tests: string | *"all" @tag(test,short=api|bench|cli|unit)
}

#Actual: #Suites & {

	if #Flags.suite == "all" {
		if #Flags.tests == "all" {
			Suites
		}

		if #Flags.tests != "all" {
			for sname, suite in Suites { 
				"\(sname)": "\(#Flags.tests)": suite[#Flags.tests]
			}
		}
	}

	if #Flags.suite != "all" {
		if #Flags.tests == "all" {
			"\(#Flags.suite)": Suites[#Flags.suite]
		}

		if #Flags.tests != "all" {
			"\(#Flags.suite)": "\(#Flags.tests)": Suites[#Flags.suite][#Flags.tests]
		}
	}

}

Sorted: {
	d1: [ for s, suite in #Actual { suite } ]
	d2: [ for suite in d1 { [for t, test in suite { test }] } ]
	f: list.FlattenN(d2, 1)
	s: list.Sort([{name: "2"}, {a: "3"}, {a: "1"}], {x: { a: string }, y: { a: string }, less: x.a < y.a}) 
	//s: list.Sort(d2, { 
		//x: {...}
		//y: {...}
		//less: x.name < y.name
	//})
}
