# Communication Subsystem (`src/comms`)

This directory is the central nervous system of the extension, managing all inter-process and intra-process communication. It provides three core services: a WebSocket client for the backend server, a global in-process event bus, and status bar management.

## Key Exports (`index.ts`)

The `index.ts` file acts as the public interface for the communication layer.

```typescript
export { extensionEmitter } from './events';
export { sendMessage } from './websocket';
export { updateStatusBar } from './statusBar';

export async function activate(context: vscode.ExtensionContext) {
  console.log("  activating comms system")

  await sbActivate(context)
  await wsActivate(context)
}

export async function deactivate() {
  await wsDeactivate()
  await sbDeactivate()
}
```

## 1. Global Event Bus (`events.ts`)

The `extensionEmitter` is the main decoupling mechanism within the extension, allowing any component to publish or subscribe to messages originating from the WebSocket.

```typescript
import * as vscode from 'vscode';

// This emitter will be used to pass messages between
// the websocket (which receives) and the sidebar (which displays).
export const extensionEmitter = new vscode.EventEmitter<any>();
```

## 2. WebSocket Client (`websocket.ts`)

This is the primary channel for talking to the external `hof agent` server.

### Message Structure
The communication protocol uses a simple `type` and `payload` object structure.
```typescript
interface Message<T> { type: string; payload: T; }
interface HelloPayload { version: string; }
interface EchoPayload { text: string; }
```

### Public API
- `sendMessage(msg)`: Sends a message to the backend server.
- **Inbound Handling**: Upon receiving a message, it is immediately parsed and forwarded to the internal bus: `extensionEmitter.fire(msg);`

```typescript
export function sendMessage<T>(msg: Message<T>) {
  if (!ws) {
    vscode.window.showErrorMessage('Server not connected.');
    return;
  }
  console.log(`[VSCODE]:`, msg);
  ws.send(JSON.stringify(msg));
}

export function sendEcho(textToSend: string) {
  if (!ws) {
    vscode.window.showErrorMessage('Server not connected.');
    return;
  }
  const payload: EchoPayload = { text: textToSend };
  const message: Message<EchoPayload> = { type: 'echo', payload };
  ws.send(JSON.stringify(message));
}
```

## 3. Handling Inbound Messages

Modules across the extension subscribe to `extensionEmitter` to react to messages from the server or other extension components. The standard pattern is to use a `switch` statement on the message `type`.

```typescript
import { extensionEmitter } from '../comms';

export function activate(context: vscode.ExtensionContext) {
  extensionEmitter.event((msg) => {
    switch (msg.type) {
      case "sync.request":
        // Handle sync request
        break;
      case "session.list.resp":
        // Handle session list response
        const sessions = msg.payload;
        break;
      default:
        // Optional: handle unknown types
        break;
    }
  });
}
```

## 4. Status Bar (`statusBar.ts`)

Manages the visual connection status in the VS Code status bar.

### Public API
```typescript
/**
 * A shared function to update the status bar text, tooltip, and icon.
 */
export function updateStatusBar(
  text: string,
  tooltip: string,
  iconName?: string // e.g., 'sync~spin', 'check', 'error'
) {
  if (!myStatusBarItem) {
    return;
  }
  myStatusBarItem.text = iconName ? `$(${iconName}) ${text}` : text;
  myStatusBarItem.tooltip = tooltip;
  myStatusBarItem.show();
}
```
