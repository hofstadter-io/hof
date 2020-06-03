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

	dir:  string
	test: string

	cover?: string
	sonar?: string
	bench?: string

	...
}

#Defaults: {
	gocli: {
		lang: "go"
		type: "testsuite"
		test: string | *"go test ./"
		...
	}
	gounit: {
		lang:  "go"
		type:  "gotest"
		test:  string | *"go test ./ -v -covermode=count"
		cover: string | *"\(test) -coverprofile cover.out && go tool cover -html=cover.out -o cover.html"
		sonar: string | *"go test ./ -v -covermode=count -coverprofile cover.out -json > tests.out"
		...
	}
}

Suites: #Suites & {

	st: {
		cli:  #Defaults.gocli  & { skip: true, dir: "lib/structural" }
		unit: #Defaults.gounit & { dir: "lib/structural" }
	}


	mod: {
		cli:  #Defaults.gocli  & { dir: "lib/mod" }
		unit: #Defaults.gounit & { skip: true, dir: "lib/mod" }
	}

}
