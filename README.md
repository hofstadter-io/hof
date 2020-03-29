# Hof - A Polyglot Code Generation Framework for Cuelang

__Hof__ is a framework for creating
code generation libraries, modules, and tools
for any technology, system, process, or platform.
Being built on [Cuelang](https://cuelang.org)
you get all the type safety and expressiveness
that comes with the language.
__Hof__ extends cuelang to make
code generation a first-class concept,
optimizing for large scale code generation.

__Hof__ also has roots in
declarative programming and DSLs,
where your write your high-level
"designs" in the "language" of a DSL.
The __hof__ tool will read your designs
and use polyglot generators to output the implementation
in one or more computer languages and technologies.
You can modify the code in the output,
update your designs, and regenrate the output.
__Hof__ allows you to work in both sides 
of the transformation, using diff3 to
merge the files together.


### Concepts in hof

Designs + Generators -> __hof__ => All the things

__Designs__ are written in DSLs using __cuelang__ syntax.
They form the "source-of-truth" for your
application, data validation, or other creation.
Designs are essentially data or configuration
with a predefined structure and their own validation from the DSL.

__Generators__ are specs and a configuration for how to generate files.
The spec creates the contract between the designs and generators, for users and writers.
The generator has input data and uses this to

__hof__ also supports modules, packages, and imports.
Quickly assemble code generators, specs, and common pieces
into advanced, cloud native applications.



### Installation

The best method is to clone and checkout the [latest release](https://github.com/hofstadter-io/hof/releases).

```
git clone https://github.com/hofstadter-io/hof
cd hof
git checkout vX.Y.Z

go mod vendor
go install
```

You can also `go get` the latest on master branch.

```
go get -u github.com/hofstadter-io/hof
```


### An example

This example is taken from the `test/templates` directory.

To run:

```
# From the hof repo root
cd test/template

# Vendor modules
hof mod vendor cue

# Generate code
hof gen

# Inspect output
cat output/*
```

Annotated file contents:

```
// Declare a package
package example

// Import modules, packages, schemas, configuration, anything cuelang
import (
  "github.com/hofstadter-io/hof/schema"
)

// Define some input data for our code generation
Input :: {
  a: "a"
  b: "b"
  N: {
    x: "x"
    y: "y"
  }
}

// Define a "HofGenName"
// The prefix is understood by hof and it will run the generator.
// Here, we are configuring an instance of a generator and giving it a name.
HofGenTest: TestGen & {
  In: {
    Val: A 
  }
  ...
}

// Define a generator
TestGen :: schema.HofGenerator & {
  Outdir: "output"

  // Generator input, this will be merged onto the Out elements
  //   (each element can define local input whih will be augmented by the generator input)
  In: {
    Val: _
    ...
  }

  // Create a list of files to generate
  Out: [
  
    // Some defaults for templates
    schema.HofGeneratorFile & {
      Template: "Val.a = '{{ .Val.a }}'\n"
      Filepath: "\(Outdir)/default.txt"
      TemplateConfig: {
        TemplateSystem: "golang"
      }
    },
    
    // Alternate delims
    schema.HofGeneratorFile & {
      Template: "Val.a = '{% .Val.a %}'\n"
      Filepath: "\(Outdir)/altdelim.txt"
      TemplateConfig: {
        AltDelims: true
        LHS2_D: "{%"
        RHS2_D: "%}"
        LHS3_D: "{%%"
        RHS3_D: "%%}"
      }
    },
    // Swap delims, using defaults delims for swap/temp
    schema.HofGeneratorFile & {
      Template: "Val.a = '{% .Val.a %}' and also this should stay {{ .Hello }}\n"
      Filepath: "\(Outdir)/swapdelim.txt"
      TemplateConfig: {
        AltDelims: true
        SwapDelims: true
        LHS2_D: "{%"
        RHS2_D: "%}"
        LHS3_D: "{%%"
        RHS3_D: "%%}"
      }
    },
    // TODO Swap delims, using custom delims for swap/temp

    // Mustache system
    schema.HofGeneratorFile & {
      Template: "Val.a = '{{ Val.a }}'\n"
      Filepath: "\(Outdir)/mustache.txt"
      TemplateConfig: {
        TemplateSystem: "raymond"
      }
    },


    // Named things
    schema.HofGeneratorFile & {
      TemplateName: "named"
      Filepath: "\(Outdir)/named-things.txt"
    },

    // File based
    schema.HofGeneratorFile & {
      TemplateName: "template-file.txt"
      Filepath: "\(Outdir)/template-file.txt"
    },

  ]

  // Templates on disk that can be referenced by name
  TemplatesDir: "templates/"
  // Partials on disk that can be referenced by name
  PartialsDir: "partials/"
  // Static files on disk to be copied into the outdir
  StaticGlobs: ["static/**"]

  // Templates in Cue that can be referenced by name
  NamedTemplates: {
    named: """
    named is '{{ .Val.a }}'
    """
  }

  // Partials in Cue that can be referenced by name
  NamedPartials: {
    named: """
    partial is '{{ .Val.a }}'
    """
  }

  // Static files in Cue to be copied into the outdir
  StaticFiles: {
    "static-cue.txt": """
    Hello, I am a static file in cue
    """
  }
}
```

