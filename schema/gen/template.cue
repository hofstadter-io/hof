package gen

// #Statics is used for static files copied over, bypassing the template engine
#Statics: {
	Globs: [...string]
	TrimPrefix: string | *""	
	OutPrefix: string | *""
}

// #Template is used for embedded or named templates
#Template: {
	Content: string
	Config?: #TemplateConfig
}

// #Templates is used for templates loaded from a filesystem
#Templates: {
	Globs: [...string]
	Config?: #TemplateConfig
}

// #TemplateConfig determines the engine and delimiters to use when rendering
// For more details see https://docs.hofstadter.io/code-generation/templates/
#TemplateConfig: {
	// The template system to use
	Engine: *"golang" | "raymond"

	// AltDelims are used when you have
	// different delims you'd like to use
	AltDelims?: #TemplateDelims

	// TmpDelims should be set to true when you have
	// output which contains the standard delims
	// which will be mistakenly processed
	// i.e.: `{{' '}}` '{{{' '}}}'
	TmpDelims: bool | *false
}

#TemplateDelims: {
  LHS2: string | *"{{"
  RHS2: string | *"}}"
  LHS3: string | *"{{{"
  RHS3: string | *"}}}"
}
