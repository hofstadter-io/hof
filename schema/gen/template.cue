package gen

#EmptyTemplates: EmptyTemplates
EmptyTemplates: {
	Templates: []
	Partials: []
	Statics: []
	...
}

#SubdirTemplates: SubdirTemplates
SubdirTemplates: {
	#subdir: string | *"."
	Templates: [{
		Globs: ["\(#subdir)/templates/**/*"]
		TrimPrefix: "\(#subdir)/templates/"
	}]
	Partials: [{
		Globs: ["\(#subdir)/partials/**/*"]
		TrimPrefix: "\(#subdir)/partials/"
	}]
	Statics: [{
		Globs: ["\(#subdir)/statics/**/*"]
		TrimPrefix: "\(#subdir)/statics/"
	}]
	...
}

#TemplateSubdirs: TemplateSubdirs
TemplateSubdirs: {
	#subdir: string | *"."
	Templates: [{
		Globs: ["./templates/\(#subdir)/**/*"]
		TrimPrefix: "./templates/\(#subdir)/"
	}]
	Partials: [{
		Globs: ["./partials/\(#subdir)/**/*"]
		TrimPrefix: "./partials/\(#subdir)/"
	}]
	Statics: [{
		Globs: ["./statics/\(#subdir)/**/*"]
		TrimPrefix: "./statics/\(#subdir)/"
	}]
	...
}

// #Statics is used for static files copied over, bypassing the template engine
#Statics: Statics
Statics: {
	Globs: [...string]
	TrimPrefix?: string
	OutPrefix?:  string
}

// #Template is used for embedded or named templates or partials
#Template: Template
Template: {
	Content: string
	Delims?: TemplateDelims
}

// #Templates is used for templates or partials loaded from a filesystem
#Templates: Template
Templates: {
	Globs: [...string]
	TrimPrefix?: string
	Delims?:     TemplateDelims

	// apply delims to a subset of templates, by glob
	DelimGlobs: [glob=string]: TemplateDelims
}

#TemplateDelims: TemplateDelims
TemplateDelims: {
	LHS: string | *""
	RHS: string | *""
}
