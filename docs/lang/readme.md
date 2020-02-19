# Hof Lang

Inspired by CueLang and Golang.

Yaml became insufficient.

Declarative programming, do want some imperative capabilities



### Files, Packages, Modules

Model after golang more than cuelang

Module is collection of packages and dependencies.

Package at directory level, must have same name for all files

File:

```bash
packege <name>

import (
  "github.com/hofstadter-io/hof/lang"
)

# Open Type
Name : (pkg.Type &)* ?{

}

# Closed Type / Value
Name :: (pkg.Type &)* ?{

}

# Template / Function ? as different from a generator template, another hame perhaps? (mould, form)
# one imperative, one declarative?
Name < args... > pkg.Type {

}

Name ( args... ) pkg.Type {

}

# Generate
Name @ (pkg.Type &)* ?{

}

# Command, as different from imperative func, as something that can be run with `hof run ....`
Name $ {

}
```


### Types, constaints, values

Should all be equivalent

Can only generate fully specified values

Types and constraints are used to enforce values


```bash
type Point struct {
  X int
  Y int
}

type PointXP Point {
  X: >0
}

type PointYP Point {
  Y: >0
}

# Omit 
type PointQ1 PointXP & PointYP 



```

Closed v Open

how to extend?

```bash
type Foo struct {
  A: string
  B: int
}

type Bar struct {
  ...Foo
  C: bool
}

type Baz Foo (&?) {
  C: bool
}
```

how to combine values?

```bash
type Cli struct {
  commands: [
    ...commonCommands.AuthCommands,
    subcmdA,
    subcmdB,
  ]
}
```

```bash
type Foo struct {
  A string
}

type Bar struct {
  B string
}

type Baz Foo & Bar {
  C string
}
```

### Scope

Public v private
 - Golang style vs data in the wild
 - only at package level?

File:
 - package
 - own imports


