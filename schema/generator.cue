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
  PackageName: string | * ""

  // Base directory of entrypoint templates to load
  TemplatesDir: string | * ""

  // Base directory of partial templatess to load
  PartialsDir: string | * ""

  // Filepath globs for static files to load
  StaticGlobs: [...string] | * [...string]

  //
  // Open for whatever else you may need
  //   often hidden fields are used
  ...
} 
