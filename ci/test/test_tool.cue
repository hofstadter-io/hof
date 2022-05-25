package test

import (
	"encoding/yaml"
	"list"
	"tool/cli"
	"tool/exec"
)

#Flags: {
	suite:  string | *"all"             @tag(suite,short=st|mod)
	tests:  string | *"all"             @tag(test,short=api|bench|cli|unit)
	subcmd: *"test" | "cover" | "sonar" @tag(subcmd,short=test|cover|sonar)
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
		text: """
		sonar: "\(#Flags.sonar)"
		\(yaml.Marshal(data))
		"""
	}
}

command: peek: {
	d1: [ for s, suite in #Actual {suite}]
	d2: [ for suite in d1 {[ for t, test in suite {test}]}]
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

command: "run-tests": {
	d1: [ for s, suite in #Actual {suite}]
	d2: [ for suite in d1 {[ for t, test in suite {test}]}]
	f: list.FlattenN(d2, 1)

	data: f

	for i, d in data {
		if d.skip == false {

			// run tests in normal mode
			if #Flags.subcmd == "test" {
				"run-\(d.name)": exec.Run & {
					cmd: ["bash", "-c", script]

					script: """
					echo "testing \(d.name)..."

					pushd ../../\(d.dir) > /dev/null
					\(d.test)
					popd > /dev/null
					"""
				}
			}

			if #Flags.subcmd == "cover" {
				if d.cover != _|_ {
					"run-\(d.name)": exec.Run & {
						cmd: ["bash", "-c", script]

						script: """
						echo "testing \(d.name)..."

						pushd ../../\(d.dir) > /dev/null
						\(d.cover)
						popd > /dev/null
						"""
					}
				}
			}

			// run tests with sonar output enabled
			if #Flags.subcmd == "sonar" {
				if d.sonar != _|_ {
					"run-\(d.name)": exec.Run & {

						cmd: ["bash", "-c", script]

						script: """
						echo "testing \(d.name)..."

						pushd ../../ >/dev/null

						mkdir -p sonar-reports/\(d.lang)/\(d.name)

						pushd \(d.dir) > /dev/null
						\(d.sonar)
						popd > /dev/null

						cp \(d.dir)/{cover,tests}.out sonar-reports/\(d.lang)/\(d.name)

						popd > /dev/null
						"""
					}
				}
			}
		}
	}
}
