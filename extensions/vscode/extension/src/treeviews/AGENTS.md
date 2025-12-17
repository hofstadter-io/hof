# Tree View Subsystem (`src/treeviews`)

This directory contains the implementations for various VS Code Tree Views, which populate the sidebars and panels defined in `package.json`. The design pattern is highly modular: each view (e.g., sessions, agents, planning) has its own file defining a `TreeDataProvider` and a custom `vscode.TreeItem` class.

All components rely on the central communication system (`../comms`) for refresh events.

## Architecture

The main entry point `index.ts` activates the views.

```typescript
// src/treeviews/index.ts
import * as vscode from 'vscode';
import { activate as sessions } from './sessions'

export async function activate(context: vscode.ExtensionContext) {
  console.log(`activating treeviews`)
	await sessions(context)
}
```

## Core Data Structures

### Example TreeItem (`sessions.ts`)

```typescript
export class Session extends vscode.TreeItem {
	constructor(
		public readonly session: any,
		public readonly label: string,
		public readonly collapsibleState: vscode.TreeItemCollapsibleState,
		public readonly command?: vscode.Command,
	) {
		super(session.sid, collapsibleState);

		this.tooltip = `${this.label}\n${this.session.sid}\n${this.session.lastUpdate}`;
		this.description = this.session.lastUpdate;
	}

	contextValue = 'session';
}
```

## Registered Tree Views

| File | View ID (package.json) | Context Value | Purpose |
| :--- | :--- | :--- | :--- |
| `sessions.ts` | `veg-sessions` | `session` | Lists active sessions. Supports commands like `veg.session.openEnviron`, `veg.session.chat`, and `veg.session.delete`. |
| `agents.ts` | `veg-agents` | `agent` | Lists available agents. |
| `planning.ts` | `veg-planning` | `task` | Displays planning/task structure. |
| `dagger.ts` | `veg-dagger` | `session` | Placeholder/duplicate of sessions. |
| `environs.ts` | `veg-environs` | `session` | Lists environments. |
| `registry.ts` | `veg-registry` | `session` | Lists registry images. |
| `session-fs.ts` | `veg-session-fs` | `agent` | File system view for sessions. |

## Command Handling

Commands are registered to handle actions on tree nodes, often firing events to the bus.

```typescript
// Example from sessions.ts
vscode.commands.registerCommand('veg.session.chat', (node: Session) => {
    vscode.commands.executeCommand('veg-chat.focus')
    const payload = { sid: node.session.sid }
    extensionEmitter.fire({ type: "chat.loadSession", payload })
    sendMessage({ type: "session.diff", payload })
});
```
