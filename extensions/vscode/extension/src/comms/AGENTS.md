# Communication Subsystem (`src/comms`)

This directory is the central nervous system of the extension, managing all inter-process and intra-process communication. It provides three core services: a WebSocket client for the backend server, a global in-process event bus, and status bar management.

## Key Exports (`index.ts`)

The `index.ts` file acts as the public interface for the communication layer.

```typescript
export { extensionEmitter } from './events';
export { sendMessage } from './websocket';
export { updateStatusBar } from './statusBar';
// ... activate functions for the whole subsystem
```

## 1. Global Event Bus (`events.ts`)

The `extensionEmitter` is the main decoupling mechanism within the extension, allowing any component to publish or subscribe to messages originating from the WebSocket.

```typescript
export const extensionEmitter = new vscode.EventEmitter<any>();
// Use: extensionEmitter.event((msg) => { ... })
// Use: extensionEmitter.fire({ type: "..." })
```

## 2. WebSocket Client (`websocket.ts`)

This is the primary channel for talking to the external `hof agent` server.

### Backend Server Management
- **Connection URL**: `ws://localhost:2257`
- **Spawning**: If a connection fails, the extension attempts to spawn the server process using:
  ```typescript
  spawn("hof", [`agent`, `--port`, `${SERVER_PORT}`], {
    detached: true,
    stdio: 'ignore',
  });
  ```
- **Status Reporting**: Connection attempts and status changes are reported via `updateStatusBar`.

### Message Structure
The communication protocol uses a simple `type` and `payload` object structure.
```typescript
interface Message<T> { type: string; payload: T; }
```

### Public API
- `sendMessage(msg)`: Sends a message to the backend server.
- **Inbound Handling**: Upon receiving a message, it is immediately parsed and forwarded to the internal bus: `extensionEmitter.fire(msg);`

## 3. Status Bar (`statusBar.ts`)

Manages the visual connection status in the VS Code status bar.

### Public API
```typescript
export function updateStatusBar(
  text: string,
  tooltip: string,
  iconName?: string // e.g., 'sync~spin', 'check', 'error'
) {
  // ... updates myStatusBarItem
}
```