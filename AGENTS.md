## hof

repo: github.com/hofstadter-io/hof
docs: hofstadter.io

### Tech Stack

Core:
- Golang
- Cuelang

Supplimentary:
- Typescript (docs, webapp, vscode)
- Dagger/Docker/Containers
- ADK
- VS Code Extension
  - webviews: pnpm, vite, react, tanstack, shadcn, tailwind

### Core Features

- **ADK + VS Code Coding Agent**
  - lib/agent/... (core runtime in Go using ADK)
  - extensions/vscode/
    - extension/... (core extension)
    - webviews/...  (webview components)
- **CUE + text/template code generation**
  - schemas/... (core schemas tied to core features)
  - lib/hof (metadata and attribute management)
  - lib/runtime (core runtime for rest of hof)
- Other features
  - ./flow (CUE base workflow engine)
  - ./formatters (containers for formatting code)
  - ./design (hof's own CUE design, just for the cli)
  - ./cmd (the generated cli code, proxy to ./lib/...)
  - ./docs (for users of hof)

### Project Organization


```sh
hof
├── .veg
│   ├── data
│   ├── embed
│   └── project
├── ci
├── cmd
│   └── hof
├── cue.mod
├── design
├── docs
├── extensions
│   └── vscode
├── flow
├── formatters
├── hack
├── images
├── lib
│   ├── agent
│   ├── chat
│   ├── config
│   ├── connector
│   ├── container
│   ├── create
│   ├── cuecmd
│   ├── cuetils
│   ├── dagger
│   ├── database
│   ├── datamodel
│   ├── datautils
│   ├── diff3
│   ├── dotpath
│   ├── extern
│   ├── fmt
│   ├── gen
│   ├── gotils
│   ├── hof
│   ├── prompt
│   ├── repos
│   ├── runtime
│   ├── singletons
│   ├── structural
│   ├── templates
│   ├── test
│   ├── tui
│   ├── types
│   ├── watch
│   └── yagu
└── schemas
    ├── chat
    ├── common
    ├── create
    ├── cue.mod
    ├── dm
    ├── gen
    ├── prompt
    └── test
```
