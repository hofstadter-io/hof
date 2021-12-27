package gen

// #Statics is used for static files copied over, bypassing the template engine
#Statics: {
	Globs: [...string]
	TrimPrefix: string | *""	
	OutPrefix: string | *""
}

// #Template is used for embedded or named templates or partials
#Template: {
	Content: string
	Delims?: #TemplateDelims
}

// #Templates is used for templates or partials loaded from a filesystem
#Templates: {
	Globs: [...string]
	TrimPrefix: string | *""	
	Delims?: #TemplateDelims
}

#TemplateDelims: {
  LHS: string | *"{{"
  RHS: string | *"}}"
}
