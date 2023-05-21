package hof

Tempate: {
	@task(hof.Template)
	$task: "hof.Template"

	name: string | *""
	data: _

	template: string
	partials: [string]: string

	out: string
}
