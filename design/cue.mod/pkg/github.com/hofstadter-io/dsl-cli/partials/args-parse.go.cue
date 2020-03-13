package partials

ArgsParse : RealArgsParse

RealArgsParse : """
{{ define "args-parse" }}
// Argument Parsing

{{ end }}
"""

OrigArgsParse : """
// Argument Parsing
{{#with . as |CMD| }}
{{#each CMD.args}}
{{#with . as |ARG|}}
	{{#if ARG.required}}
	if {{@index}} >= len(args) {
		fmt.Println("missing required argument: '{{ARG.name}}'\n")
		cmd.Usage()
		os.Exit(1)
	}
	{{/if}}

	var {{camel ARG.name}} {{> go-type.go ARG.type}}
	{{#if ARG.default}}
		{{#if (eq ARG.type "string")}}
		{{camel ARG.name}} = "{{ARG.default}}"
		{{else}}
		{{camel ARG.name}} = {{ARG.default}}
		{{/if}}
	{{/if}}

	if {{@index}} < len(args) {
	{{#if ARG.rest}}
		{{#if (eq ARG.type "array:string")}}
			{{camel ARG.name}} = args[{{@index}}:]
		{{else}}
		{{/if}}
	{{else if (eq ARG.type "string")}}
			{{camel ARG.name}} = args[{{@index}}]
	{{else}}
			{{camel ARG.name}}Arg := args[{{@index}}]
			var err error
			{{> common/golang/parse/builtin.go IN_NAME=(concat2 (camel ARG.name) "Arg") OUT_NAME=(camel ARG.name) TYP=ARG.type}}
			if err != nil {
				fmt.Printf("argument of wrong type. expected: '{{ARG.type}}' got error: %v", err)
				cmd.Usage()
				os.Exit(1)
			}
	{{/if}}
	}
{{/with}}
{{/each}}
{{/with}}

"""
