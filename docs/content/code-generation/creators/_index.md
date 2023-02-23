---
title: "Creators"
linkTitle: "Creators"
weight: 61
---

{{<lead>}}
`hof create` is a command for bootstrapping
from user input from flags, files, or a prompt.
They are also flexible enough for any git repository
to provide the creator experience.
Create one-liners for your users so they can
get started with your projects quickly and easily.
{{</lead>}}

A `creator` is a hof generator with a `Create` field,
run as a one-time code generation.
However, a 
Most hof code generators will also provide a creator
to help their users get started easier.

- The main part is the `prompt` where you can ask the user for inputs
- These then form the input to the file bootstrapping process (one-time hof code gen)
- A `hof generator` will use `hof create` to create the initial CUE files for using their generator
- You can add `hof create` to any repository to generate any files

Users then use `hof create` bootstrap new applications, modules, or configuration
like files needed for CI/CD, deployment, security, or other system.

{{<childpages>}}

<div style="margin-bottom: 13rem;"></div>

{{< youtube TfaEV37C6IE >}}
