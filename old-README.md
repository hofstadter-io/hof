# Hof - Polyglot Code Generation Framework

Hof is a Polyglot Code Generate Framework.
You write designs in Cue, very similar to yaml or json,
and Hofstadter validates and feeds these into
generators to create directories and files.
These directories and files can be what ever you choose,
and language, technology or combination.

Under the hood, Hofstadter is using a templating engine.
Cuelang is used to provide type safety with incredible expressiveness.
Modules can be used like other languages
to share and import functionality from small to entire applications and systems.
You can write custom code in the generated output while
still modifying design and regenerating code.
Hofstadter manages the diff process to ensure your work is not disrupted.

The Hofstadter Framework provides conventions and a specification
for writing generators. It wraps Cuelang and the templating engine
to enable you to generate code at scale.
Feed directories of data through the same generator,
the same data through multiple generators,
or any combination.

Hofstadter extends the flexibility and safety of Cue to code generation.


### [Documentation](https://hofstadter.io/docs)

Alternate description:

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

Prebuilt binaries are available in the [Releases section](https://github.com/hofstadter-io/hof/releases/).
There are also [docker images available](https://hub.docker.com/repository/docker/hofstadter/hof).

To install from source:

```
git clone https://github.com/hofstadter-io/hof
cd hof
git checkout vX.Y.Z

go install
```

You can also `go get` the latest on master branch.

```
go get -u github.com/hofstadter-io/hof
```


### An example

This example is taken from the `test/templates` directory.

You can find the schemas for `HofGenerator` and `HofGeneratorFile`
in the [schemas directory](./schema).

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

Output:

```
$ tree output
output/
├── altdelim.txt
├── default.txt
├── mustache.txt
├── named-things.txt
├── static-cue.txt
├── static-file.txt
├── swapdelim.txt
└── template-file.txt

$ find output -type f -exec echo "{}" \; -exec cat "{}" \; -exec echo \;
output/swapdelim.txt
Val.a = 'a' and also this should stay {{ .Hello }}

output/named-things.txt
named is 'a'
output/template-file.txt
Hi, I'm a template file on disk 'a'
 ... and I'm a partial Hi, I'm a partial file on disk 'b'


---
HellWorld
---

output/static-file.txt
Hello, I am a static file from the filesystem

output/mustache.txt
Val.a = 'a'

output/default.txt
Val.a = 'a'

output/static-cue.txt
Hello, I am a static file in cue
output/altdelim.txt
Val.a = 'a'


```

### Directory structure

`schema` holds the Cue code which implements the very core and foundational concepts in hofland.

`hof` main development files
- `lib` houses much of the hand written code, libraries, and core logic 
- `docs` holds some project oriented development. There is a full site dedicated to user facing docs.
- `test` contains many cases, files, scenarios, and drivers for testing.

`hof` generated files (or self-referential development;)
- `design` contains the Cue code for the capabilities generated here
- `ci` files for cloud builds with mainly testing
- `cmd` files for the cli structure and helpers
- `gen` files implementing the other capabilities here

### Modules and Examples

Projects:

- This project uses itself to generate various pieces like the cli structure and the release process.

Modules:

- [hofmod-model](https://github.com/hofstadter-io/hofmod-model) - A module for representing common types and their relations.
- [hofmod-cli](https://github.com/hofstadter-io/hofmod-cli) - Create CLI infrastructure based on the Golang Cobra library.
- [hofmod-releaser](https://github.com/hofstadter-io/hofmod-releaser) - Release code or binaries to GitHub and Docker with minimal configuration. Based on [GoReleaser](https://goreleaser.com/).
- [hofmod-config](https://github.com/hofstadter-io/hofmod-config) - Cloud native config and secret files using the Golang Viper library and adding dynamic reload in Kubernetes.
- [hofmod-rest](https://github.com/hofstadter-io/hofmod-rest) - Generate Golang REST servers that are ready for production. This makes use of many of the other modules here.
- [hofmod-hugo](https://github.com/hofstadter-io/hofmod-hugo) - Create documenation sites with [Hugo](https://gohugo.io) and [Docsy](https://docsy.dev)


