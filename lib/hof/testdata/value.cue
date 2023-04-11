package hack

val: {
	foo: "bar"

	$hof: {
		metadata: {
			labels: {
				app: "val-app"
				env: "dev"
			}
		}
		datamodel: {
			root:    true
			history: true
		}
	}
}
