package test

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
			cmd: "go test ./ -v -covermode=count"
			sonar: "go test ./ -v -covermode=count -coverprofile cover.out -json > tests.out"
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
			cmd: "go test ./ -v -covermode=count"
			sonar: "go test ./ -v -covermode=count -coverprofile cover.out -json > tests.out"
		}
	}

}
