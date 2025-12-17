# Source Code Root (`src/`)

This directory contains the main implementation for the Veg VS Code Extension, organized into modular subsystems that handle communication, state synchronization, service provision, and user interface components.

## Architecture

The extension follows a classic layered structure, orchestrated by the main entry point, `extension.ts`.

### Activation Order (`extension.ts`)

The `activate` function establishes the core dependencies in a specific order:

```typescript
import * as comms from './comms'
import * as sync from './sync'
import * as treeviews from './treeviews';
import * as webviews from './webviews';
import * as filesys from './services/filesystemProvider'

export async function activate(context: vscode.ExtensionContext) {
	// 1. Establish communication channels (websocket, internal emitter)
	await comms.activate(context)
	// 2. Register the core service (veg:// filesystem)
	await filesys.activate(context)
	// 3. Start background monitoring (env, window, terminals, workspace)
	sync.activate(context)
	// 4. Register UI components
	await webviews.activate(context)
	await treeviews.activate(context)
	// 5. Trigger initial state synchronization with the backend
	comms.extensionEmitter.fire({ type: "requestSync" })
	comms.sendMessage({ type: "requestSync", payload: {} })
}
```

## Subsystem Index

| Directory | Purpose | Link |
| :--- | :--- | :--- |
| `comms` | Manages WebSocket connection to the backend and the internal event bus. | [./comms/AGENTS.md](./comms/AGENTS.md) |
| `services` | Provides the `veg://` virtual filesystem service and associated commands. | [./services/AGENTS.md](./services/AGENTS.md) |
| `sync` | Monitors and broadcasts VS Code state (workspace, terminals, environment) to the server. | [./sync/AGENTS.md](./sync/AGENTS.md) |
| `treeviews` | Registers and manages all `vscode.TreeDataProvider` implementations (sessions, agents, planning). | [./treeviews/AGENTS.md](./treeviews/AGENTS.md) |
| `webviews` | Implements the framework for all sidebars and panels (`chat`, `debug`, `manage`). | [./webviews/AGENTS.md](./webviews/AGENTS.md) |

## Utilities and Tests

- **`other/token-count.ts`**: Utilities for counting tokens (likely for AI context management).
- **`test/extension.test.ts`**: Standard VS Code extension integration tests.
