# Synchronization Subsystem (`src/sync`)

This directory is the core state synchronization layer, responsible for observing the VS Code environment and broadcasting state changes to the backend server (via the `../comms` system).

All sync modules are activated by `src/sync/index.ts` and triggered primarily by `requestSync` events.

## Core Components

| File | Purpose | Data Broadcasted | Key Features |
| :--- | :--- | :--- | :--- |
| `env.ts` | **Environment** | `sid`, `machineId`, `sessionId`, `workspaceDir`, clipboard content. | Provides essential metadata about the VS Code instance and context. |
| `window.ts` | **Window State** | Tab group structure, active/dirty tabs, visible text editors, window state. | Captures the UI layout and editing activity. |
| `workspace.ts` | **Workspace** | Workspace name, trust status, folders, open text documents. | Captures the project-level context. |
| `terminals.ts` | **Terminal Activity** | Command line, current working directory, output, exit code for each execution. | Tracks terminal commands using VS Code Shell Integration. Also handles opening a Dagger-based terminal via the `session.term.open` event. |

## Legacy/Unused Components

- **`fs.ts`**: Contains a `VegContentProvider` implementation similar to `src/services/filesystemProvider.ts`. It is **NOT** currently imported or activated by `index.ts`. It appears to be an older or alternative implementation of the virtual filesystem logic.

## Terminal Tracking Details (`terminals.ts`)

The terminal synchronization is robust, capturing full execution history.

```typescript
class Exec {
	start: vscode.TerminalShellExecutionStartEvent | null = null
	output: string = ""
	end: vscode.TerminalShellExecutionEndEvent | null = null
}

class Terminal {
	termIndex: number = 0;
	terminal: vscode.Terminal | null = null;
	history: Exec[] = [];	
}

class TerminalPayload {
	id: number = -1;
	name?: string;
	history?: HistoryPayload[];
}

class HistoryPayload {
	cmd?: any;
	cwd?: string;
	out?: string;
	exit?: number;
}
```

It listens to:
- `onDidStartTerminalShellExecution`
- `onDidEndTerminalShellExecution`
- `read()` from the execution stream to capture output.

## Environment Synchronization (`env.ts`)

Broadcasts information about the environment.

```typescript
async function broadcastEnv(context: vscode.ExtensionContext) {
	const wsF = vscode.workspace.workspaceFolders
	var wDir: string | undefined
	if (wsF && wsF.length > 0) {
		wDir = wsF[0].uri.path	
	}
	const sid = context.workspaceState.get("sid")

	const msg = {
		type: "env.info.resp",
		payload: {
			sid,
			machineId: vscode.env.machineId,
			vscodeSid: vscode.env.sessionId,
			remoteName: vscode.env.remoteName,
			user: "verdverm",
		  workspaceDir: wDir,
			clipboard: await vscode.env.clipboard.readText(),
		}
	}
	extensionEmitter.fire(msg);
	sendMessage(msg)
}
```
