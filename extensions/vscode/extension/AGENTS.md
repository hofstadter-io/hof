# Veg VS Code Extension

The Veg extension provides an interactive, Dagger-powered workspace for software engineering tasks. It integrates remote environment state, chat-based interaction, and a custom virtual filesystem directly into the VS Code environment.

## Architecture and Codebase

The entire functional codebase resides in the [src/ directory](./src/AGENTS.md), which is organized into modular subsystems for communication, state synchronization, service provision, and user interface components.

## Key Metadata (`package.json`)

| Field | Value | Purpose |
| :--- | :--- | :--- |
| `name` | `veg` | Project name. |
| `displayName` | `Veggie` | Name displayed in the VS Code marketplace/UI. |
| `main` | `./out/extension.js` | Compiled entry point for the extension host. |
| `activationEvents` | `onStartupFinished` | The extension is loaded when VS Code finishes starting up. |

### UI Contributions (Views)

The extension registers three main view containers:

1.  **`veg-chat-sidebar` (Activity Bar)**: Contains the `veg-chat` webview.
2.  **`veg-manage-sidebar` (Activity Bar)**: Contains the `veg-manage` webview.
3.  **`veg-debug-panel` (Panel)**: Contains the `veg-debug` webview and the `veg-sessions` tree view.

### Dependencies

The extension relies on two major external dependencies:
- `ws`: For the WebSocket connection to the local backend server (`hof agent`).
- `jsonc-parser`: For handling configuration and parsing JSON with comments.

### Virtual Filesystem and URI Schemes

The extension manages a virtual filesystem under the `veg://` scheme, which maps to the backend's `oci://` environment URIs.

- **`veg://`**: Used internally by VS Code. The authority and path segments typically encode the environment ID and version (e.g., `veg://host/envId:ver/path/to/file`).
- **`oci://`**: Used by the backend services (Dagger). The extension translates `veg://` URIs to `oci://` URIs before making API calls.
- **Transformation**: Logic in `src/services/utils.ts` (`vsUriToVeg`) handles this conversion, ensuring the backend always receives `oci://` URIs.
- **Session IDs**: The frontend explicitly resolves the active Session ID for a given URI and passes it as a separate `sid` field in API requests. This ensures that operations are session-aware (updating state) when appropriate, while allowing stateless access to arbitrary OCI images when no session is involved.

### Commands and Keybindings

Critical commands are registered for user interaction:
- **`veg.connect`**: Triggers the connection/re-connection process.
- **`veg-chat.focus`**: Keybindings for quickly opening the chat panel (`ctrl+g`, `alt+g g`).
- **`veg.explorer.*`**: Commands bound to the file explorer context menu for managing virtual environment files (e.g., `showDiff`, `mergeDiff`, `openEnviron`).
- **`veg.session.*`**: Commands bound to the session tree view for managing individual sessions (e.g., `chat`, `fork`, `delete`).

## Source Code Index

All core logic is contained within the following directories:

- [**src/**](./src/AGENTS.md)
