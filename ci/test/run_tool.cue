package test

import (
	"encoding/yaml"
	"list"
	"tool/cli"
	"tool/exec"
)

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

command: info: {
	data: #Actual

	print: cli.Print & {
		text: yaml.Marshal(data)
	}
}

command: peek: {
	d1: [ for s, suite in #Actual { suite } ]
	d2: [ for suite in d1 { [for t, test in suite { test }] } ]
	data: list.FlattenN(d2, 1)

	for i, d in data {
		"print-\(i)": cli.Print & {
			text: """
			---
			\(yaml.Marshal(d))
			"""
		}
	}
}

command: run: {
	d1: [ for s, suite in #Actual { suite } ]
	d2: [ for suite in d1 { [for t, test in suite { test }] } ]
	f: list.FlattenN(d2, 1)

	data: f

	for i, d in data {
		"run-\(i)": exec.Run & {
			script: """
			echo "testing \(d.name)..."
			"""
			cmd: ["bash", "-c", script]
		}
	}
}

