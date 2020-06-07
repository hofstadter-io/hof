package hof

generated: _ @test(suite)
generated: {
	cli: _ @test(scenario)
	cli: {
		// TODO before / after
		cmds: _ @test(script)
		cmds: {
			dir: "cmd/hof/cmd"
			scripts: [ "**/*.txt" ]
			env: {
				HELLO: "WORLD"
				FOO: "$FOO"
			}
		}
	}
}

human: _ @test(suite)
human: {

}
