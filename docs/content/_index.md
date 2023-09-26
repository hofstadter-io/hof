---
title: Hofstadter Documentation
description: "Hofstadter Documentation"
---

{{<lead>}}
Welcome to the documentation site.
{{</lead>}}

# `hof` is CUE powered tools and capabilities

`hof` helps you use CUE to enhance your work with
data models, schemas, code generation, task workflows, and more.

__`hof` is a CLI tool you *add* to your workflow.__

- Augment data, config, and schemas with CUE to improve consistency, gain confidence, and unlock new possibilities
- Generate multiple Yaml and JSON files in one-shot, from a CUE based source of truth
- Generate much of the application code, work directly in the output, regenerate without loss
- Explore data or config with the TUI, work with it using CUE in live-reload mode
- Run workflows with automatic task dependency inference, the right tasks are run in the right order

|  | |
|:---                    |:-- |
| __code generation__    | Data + templates = _ (anything), technology agnostic |
| __data model__         | Define and manage data models, life-cycle features |
| __task engine__        | Extensible task and DAG workflow engine |
| __CUE cmds__           | Core def, eval, export, and vet commands |
| __creators__           | bootstraping and starter kits from any repo |
| __modules__            | CUE module dependency management |
| __tui__                | A terminal interface to Hof and CUE |
| __chat__               | Combine LLM and Hof code gen for better, scalable results |

<br>

#### There are two modes to use `hof`

1. creating applications (green boxes)
1. building reusable modules (blue boxes)

Like most languages and frameworks, there are two types of users.
Most users will build applications for some purpose, using libraries written by others.
A smaller number will build reusable modules, like the packages and libraries you use today.
`hof` has the same for same relationship for code generators modules.
All modules exist outside of the `hof` tool and just need to be a git repository.


## Designed to augment your workflows

__`hof` is a CLI tool you will add to your workflows.__
We know developers have their own preferences
for tools, languages, and platforms.
`hof` can work with any of them.
You will typically use `hof` at development time,
committing the generated code to git.

__`hof` is technology agnostic.__
You can generate code for any language or technology,
and more often than not you will generate several together.
From your data models, the source of truth,
`hof` can generate consistent code across the stack.

__`hof` captures common patterns and boilerplate.__
Through the templates and code generation modules,
so we can remove much of the repetitive tasks and coding effort.
Updates to the data model can be replicated instantly through the stack.

__`hof` modules span technologies.__
With composable modules, we can create full-stack applications
and the infrastructure to run them by importing from the ecosystem.
Logical application features can be composed
as bigger building blocks from any language, framework, or tool.

__`hof` continues to work as your model evolves.__
Rather than a one-time bootstrapping at the beginning of development,
you can update your designs or data model and regenerate code.
Think of code generated with `hof` as living boilerplate or scaffolding.
You can also add custom code directly in the output
and `hof` will ensure it stays as you regenerate your application.


# We call this High Code development.

{{<lead>}}
Creating code with higher levels of design, reuse, and implementation
{{</lead>}}



## What can you do with `hof`?

<br>

##### Generate anything

Applications all start as files and `hof` generates directories of files.
You can generate the source files, configuration, deployment, and CI files needed.
If it's made of files, you can generate it with `hof`.

##### Consolidate the data model

The same data model appears at each level of the tech stack.
You should only have to write it down once, as a _single-source of truth_.
More than just the shape, this should also include the rules.

##### Capture common code and application patterns

Whether it is writing api handlers, CRUD, client libraries, or data validation,
there are many patterns per data model.
There are also application wide patterns.
When starting server setup like logging and wiring up the router.

##### Manage model and application versions.

Data models evolve with an application and need management.
From updating the code and databased to deployment updates and supporting
older clients, you can have multiple versions being referenced.
You latest backend will need to handle many previous versions.

##### Work directly in the (re)generated code

With `hof` you write custom code directly in the generated output,
where it naturally belongs. Your final code should look the same.
When you change your data model or designs, `hof` uses diff3
to ensure your code is left in place and 

##### Share and control modules with dependency management

Sharing models and code generation is core to `hof`
and central to solving problems of interoperability between
different teams and services.
Both design and generators are managed with versions
and dependency management.

##### Apply fleet wide fixes and updates

Deploying shared security and bug fixes across many applications should be easier.
This should apply equally for improvements in our code patterns and practices.

##### Extensible generators and models

Both generators and models can be combined through dependencies and imports.
You can extend, override, or otherwise customize as you need.
They are separate from the tool so there is no need to change `hof` 
to enable new technologies or patterns.


