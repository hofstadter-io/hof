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
- Agent Development Kit (ADK-Go)
- VS Code Extension
  - webviews: pnpm, vite, react, tanstack, shadcn, tailwind


### Core Features

We are currently developing a vscode extension for a custom coding agent setup, a copilot alternative and then some.

- [lib/agent](./lib/agent/AGENTS.md) is the backend server
- [extensions/vscode/extension](./extensions/vscode/extension/AGENTS.md) is the extension core
- [extensions/vscode/webviews/chat](./extensions/vscode/webviews/chat/AGENTS.md) is the chat interface

This is all we are working on currently.


### Project Organization

- **ADK + VS Code Coding Agent**
  - lib/agent/... (core runtime in Go using ADK)
  - lib/agent/config/... (CUE-based configuration for agents, tools, and environments)
  - extensions/vscode/
    - extension/... (core extension)
    - webviews/...  (webview components)
- **CUE + text/template code generation**
  - schemas/... (core schemas tied to core features)
  - lib/hof (metadata and attribute management)
  - lib/runtime (core runtime for rest of hof)
  - lib/gen (core generator types and logic)
- Other features
  - ./flow (CUE base workflow engine)
  - ./formatters (containers for formatting code)
  - ./design (hof's own CUE design, just for the cli)
  - ./cmd (the generated cli code, proxy to ./lib/...)
  - ./docs (for users of hof)


### Build, Test, Validate

Right now, you should only run one of the following commands.
If more are needed, the user will instruct you.

```sh
# Build, this is the main way we validate agent/ai changes work on a basic level
go install ./cmd/hof

# basic check
hof version
```

Do not run any other hof commands or go test UNLESS specifically instructed to by the user.