---
title: "CUE"
description: "A declarative language for configuration and code generation."
brief: "A declarative language for configuration and code generation."

weight: 5
---

{{<lead>}}
`hof` leverages the power of CUE (Configure, Unify, Execute) to define inputs, schemas, and generators.
CUE is a language specifically designed for configuration and is an integral component of the `hof` developer experience.
{{</lead>}}

Two key considerations drove the choice to use CUE:

1. __Declarative Code Generation__<br>
   `hof` uses code generators defined in CUE, combining data and templates to generate files.
   With CUE as the source of truth, we gain consistency for inputs and mappings, allowing
    composable generators to generate code across languages and technologies.
    `hof` allows you to define the declarative schemas to code generators and package them into reusable modules.

2. __Purpose-Built for Large-Scale Configuration__<br>
   CUE has a strong foundation and robust tooling for managing configuration across multiple files, packages, and modules.
   CUE's logical nature ensures that the configuration is both efficient and correct.
    Since `hof`'s inputs are typically declarative configurations, CUE is a natural fit throughout the process.

Other benefits of CUE include:

- Configuration-like syntax with a whole, logical language.
- Built on sound principles and rich heritage.
- Declarative and logical, providing confidence in the validity and consistency of the configuration.
- Purpose-built language and tooling for large-scale configuration.
- Dependency management for data, schemas, configuration, and code generation.

For those unfamiliar with CUE, Bitfield Consulting has a great
[Introduction to CUE from a JSON perspective](https://bitfieldconsulting.com/golang/cuelang-exciting).
It is a quick read and provides an excellent foundation for using CUE with `hof`.

To further your CUE knowledge, be sure to check out these resources:

- [CUE documentation](https://cuelang.org) (from the CUE Team)
- [Cuetorials](https://cuetorials.com) (by Hofstadter)
