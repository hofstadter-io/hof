// Argument Parsing
{{#with . as |Cmd| }}
{{#each Cmd.args}}
{{#with . as |arg|}}
	{{#if arg.required}}
	if {{@index}} >= len(args) {
		fmt.Println("missing required argument: '{{arg.name}}'\n")
		cmd.Usage()
		os.Exit(1)
	}
	{{/if}}

	var {{camel arg.name}} {{> go-type.go arg.type}}
	{{#if arg.default}}
		{{#if (eq arg.type "string")}}
		{{camel arg.name}} = "{{arg.default}}"
		{{else}}
		{{camel arg.name}} = {{arg.default}}
		{{/if}}
	{{/if}}

	if {{@index}} < len(args) {
	{{#if arg.rest}}
		{{#if (eq arg.type "array:string")}}
			{{camel arg.name}} = args[{{@index}}:]
		{{else}}
		{{/if}}
	{{else if (eq arg.type "string")}}
			{{camel arg.name}} = args[{{@index}}]
	{{else}}
			{{camel arg.name}}Arg := args[{{@index}}]
			var err error
			{{> common/golang/parse/builtin.go IN_NAME=(concat2 (camel arg.name) "Arg") OUT_NAME=(camel arg.name) TYP=arg.type}}
			if err != nil {
				fmt.Printf("argument of wrong type. expected: '{{arg.type}}' got error: %v", err)
				cmd.Usage()
				os.Exit(1)
			}
	{{/if}}
	}
{{/with}}
{{/each}}
{{/with}}

