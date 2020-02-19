{{#if (eq TYP "string")}}
	{{OUT_NAME}} = {{IN_NAME}}
{{else if (eq TYP "int")}}
	{{OUT_NAME}}_int64, err := strconv.ParseInt({{IN_NAME}}, 10, 64)
	{{OUT_NAME}} = int({{OUT_NAME}}_int64)
{{else if (eq TYP "uint")}}
	{{OUT_NAME}}_uint64, err := strconv.ParseUint({{IN_NAME}}, 10, 64)
	{{OUT_NAME}} = int({{OUT_NAME}}_uint64)
{{else}}
// UNKNOWN TYPE: '{{TYP}}'
{{/if}}

