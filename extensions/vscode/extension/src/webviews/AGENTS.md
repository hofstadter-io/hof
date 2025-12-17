# Webview Subsystem (`src/webviews`)

This directory is responsible for setting up and managing all VS Code Webview views (sidebars and panels) used by the extension. It acts as the bridge between the extension's backend logic (running in the extension host process) and the frontend UI (running in a separate, isolated browser environment).

All webviews utilize the central communication system (`../comms`) for sending/receiving data.

## Core Components

The architecture is based on a reusable `WebviewProvider` and a utility `loader` for asset management.

### `WebviewProvider` (`provider.ts`)

This class implements the `vscode.WebviewViewProvider` interface and serves as the backbone for all registered webview views.

```typescript
export class WebviewProvider implements vscode.WebviewViewProvider {
  private _view?: vscode.WebviewView;

  constructor(
    private readonly _context: vscode.ExtensionContext,
    public readonly _name: string,
    private _onMessage: (data: any) => void
  ) {}

  public async resolveWebviewView(
    webviewView: vscode.WebviewView,
    context: vscode.WebviewViewResolveContext,
    _token: vscode.CancellationToken
  ) {
    // ... setup
    // 1. Listen for messages from the webview (UI -> Extension)
    webviewView.webview.onDidReceiveMessage(this._onMessage);

    // 2. Listen for messages from the websocket (Extension -> UI)
    extensionEmitter.event((e) => {
        this._view?.webview.postMessage(e);
    });
  }
}
```

Key responsibilities:
1.  **HTML Loading**: Uses `getHtmlForWebview` to load the built HTML assets.
2.  **Webview -> Extension**: Registers the provided `_onMessage` callback to handle messages from the webview UI.
3.  **Extension -> Webview**: Subscribes to the global `extensionEmitter` event to post messages to the webview UI.

### HTML Loader (`loader.ts`)

The `getHtmlForWebview` utility is critical for ensuring the webview can load its assets securely. It reads the pre-built `index.html` file, uses `webview.asWebviewUri` to resolve all relative `src` and `href` paths to VS Code resource URIs, and applies a `nonce` for Script-Src to comply with the Content Security Policy (CSP).

### Registered Webviews

| File | View ID (package.json) | Purpose |
| :--- | :--- | :--- |
| `chat.ts` | `veg-chat` (Sidebar) | The main chat interface. Listens for `chat.loadSession` to store the active session ID (`sid`) in `context.workspaceState`. |
| `debug.ts` | `veg-debug` (Panel) | A panel for debugging/synchronization. Provides the `veg.debug.requestSync` command, which uses the stored `sid` to request a sync from the server. Clears `sid` on `session.delete`. |
| `manage.ts` | `veg-manage` (Sidebar) | Currently a placeholder. |

Messages sent from the Chat webview (`chat.ts`) are proxied through `onMessage` to both the global `extensionEmitter` and the `websocket` via `sendMessage`.

## Public Interfaces (`index.ts`)

```typescript
import * as vscode from 'vscode';

import { activate as manage } from './manage'
import { activate as debug } from './debug'
import { activate as chat } from './chat'

export function activate(context: vscode.ExtensionContext) {
  console.log(`activating webviews`)

  debug(context)
  manage(context)
  chat(context)

}
```
