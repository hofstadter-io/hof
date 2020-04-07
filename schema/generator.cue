package schema


// Definition for a generator
HofGenerator :: {
  // Base directory for the output
  Outdir: string | *"./"

  // "Global" input, merged with out replacing onto the files
  In: { ... } | * {...}

  // The list fo files for hof to generate
  Out: [...HofGeneratorFile] | *[...]

  // Subgenerators for composition
  Generators: [...HofGenerator] | *{...}

  // The following will be automatically added to the template context
  // under its name for reference in GenFiles  and partials in templates
  NamedTemplates: { [Name=string]: string }
  NamedPartials:  { [Name=string]: string }
  // Static files are available for pure cue generators that want to have static files
  // These should be named by their filepath, but be the content of the file
  StaticFiles: { [Name=string]:  string }

  //
  // For file based generators
  //   files here will be automatically added to the template context
  //   under its filepath for reference in GenFiles and partials in templates

  // Used for indexing into the vendor directory...
  PackageName?: string

  // Base directory of entrypoint templates to load
  TemplatesDir: string | * "templates"

  // Base directory of partial templatess to load
  PartialsDir: string | * "partials"

  // Filepath globs for static files to load
  StaticGlobs: [...string] | *[]

  TemplateConfig: {
    // Include Common attributes
    // System params
    TemplateSystem: *"golang" | "raymond"

    //
    // Template delimiters
    //
    //   these are for advanced usage, you shouldn't have to modify them normally

    // Alt and Swap Delims,
    //   becuase the defaulttemplate systems use `{{` and `}}`
    //   and you may choose to use other delimiters, but the lookup system is still based on the template system
    //   and if you want to preserve those, we need three sets of delimiters
    AltDelims:  bool | *false
    SwapDelims: bool | *false

    // The default delimiters
    // You should change these when using alternative style like jinjas {% ... %}
    // They also need to be different when using the swap system
    LHS_D: LHS2_D
    RHS_D: RHS2_D
    LHS2_D: string | *"{{"
    RHS2_D: string | *"}}"
    LHS3_D: string | *"{{{"
    RHS3_D: string | *"}}}"

    // These are the same as the default becuase
    // the current template systems require these.
    //   So these should really never change or be overriden until there is a new template system
    //     supporting setting the delimiters dynamicalldelimiters dynamicallyy
    LHS_S: LHS2_S
    RHS_S: RHS2_S
    LHS2_S: string | *"{{"
    RHS2_S: string | *"}}"
    LHS3_S: string | *"{{{"
    RHS3_S: string | *"}}}"

    // The temporary delims to replace swap with while also swapping
    // the defaults you set to the swap that is required by the current templet systems
    // You need this when you are double templating a file and the top-level system is not the default
    LHS_T: LHS2_T
    RHS_T: RHS2_T
    LHS2_T: string | *"#_hof_l2_#"
    RHS2_T: string | *"#_hof_r2_#"
    LHS3_T: string | *"#_hof_l3_#"
    RHS3_T: string | *"#_hof_r3_#"
  }

  //
  // Open for whatever else you may need as a generator writer
  ...
} 
