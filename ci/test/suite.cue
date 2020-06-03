package test

#Suites: [S=string]: [T=string]: #Tester & {
	name: "\(S)/\(T)"
	...
}

#Tester: {
	name: string
	lang: string
	type: string

	skip: bool | *false

	dir: string
	cmd: string

	sonar?: string
	bench?: string

	...
}

Suites: #Suites & {

	st: {
		cli: {
			skip: true
			lang: "go"
			type: "testsuite"
			dir: "lib/structural"
			cmd: "go test ./"
		}

		unit: {
			lang: "go"
			type: "gotest"
			dir: "lib/structural"
			cmd: "go test ./ -v -covermode=count"
			sonar: "go test ./ -v -covermode=count -coverprofile cover.out -json > tests.out"
		}
	}


	mod: {
		cli: {
			lang: "go"
			type: "testsuite"
			dir: "lib/mod"
			cmd: "go test ./"
		}

		unit: {
			skip: true
			lang: "go"
			type: "gotest"
			dir: "lib/mod"
			cmd: "go test ./ -v -covermode=count"
			sonar: "go test ./ -v -covermode=count -coverprofile cover.out -json > tests.out"
		}
	}

}
