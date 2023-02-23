---
title: "How hof works"
linkTitle: "How hof works"
weight: 15
draft: true
description: >
  How code generation works.
---

Walkthrough of the processing steps.

- Load entrypoints
- Search for `@gen()` attributes
- Parse generators
- Render files
- Merge output

Per generator:

- find any subgenerators
- process files

Per file:

- Render
- Compare to existing
- if different
- diff3
- write shadow
- write merged
