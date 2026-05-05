# Veggie Chat Webview - Agent Guide

This document provides a high-level overview of the `veg-chat-webview` codebase to assist agents in understanding the project structure, architecture, and common patterns.

## Project Overview

This project is a **React-based VS Code Webview** that serves as the UI for the Veggie AI coding assistant. It communicates with a host VS Code extension backend via message passing.

**Tech Stack:**
- **Package Manager:** pnpm (workspace) 
- **Framework:** React (Vite)
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **UI Library:** Shadcn/UI (Radix primitives), Lucide React (icons)
- **State Management:** React Context (`ChatProvider`)

## Testing

IMPORTANT: you cannot run these commands yet, always ask the user to build, test, and validate.

```sh
# build the extension webapp
pnpm build
```

## Architecture

### 1. Communication Layer
The webview acts as a "dumb" terminal that renders state provided by the VS Code extension.
- **`src/vscodeApi.ts`**: wrapper around the native VS Code Webview API (`acquireVsCodeApi`).
- **`src/hooks/useChat.tsx`**: The central nervous system.
    - Listens for messages from the extension (`vscodeApi.onMessage`).
    - Sends user actions back to the extension (`vscodeApi.postMessage`).
    - Manages the global `ChatContext`.

### 2. State Management (`ChatProvider`)
The application state is centralized in `src/hooks/useChat.tsx`.
- **`session`**: The current chat session object (ID, events, metadata).
- **`sid`**: Session ID.
- **`chatState`**: Configuration, environment info, available models/agents.
- **`usage`**: Token usage statistics.
- **`diff`**: Information about pending file changes.
- **`pos`**: Time-travel cursor position (for viewing history).

### 3. Component Tree
- **`App.tsx`**: Layout root. Uses `react-resizable-panels` to split the view into:
    - Top: Chat History (`StickToBottom`).
    - Bottom: User Input.
- **`components/Events/`**: Renders the chat stream.
    - `index.tsx`: Main list of events. Handles event coalescing (merging function calls with results).
    - `Messages.tsx`: Renders individual `UserMessage` or `ModelMessage`.
    - `ToolCall.tsx`: Renders tool executions (`exec`, `fs_read`, etc.).
    - `Details.tsx`: JSON inspector and metadata for events.
- **`components/UserInput/`**: The input area.
    - `index.tsx`: Container for the editor and send controls. Handles special commands (e.g., `@agent`, `/command`).
    - `editor.tsx`: Rich text editor (Tiptap based).
- **`components/Header.tsx`**: Session title, usage stats, and session management controls (delete, new).

## Data Flow

**Sending a Message:**
1. User types in `UserInput`.
2. `handleSend` in `useChat` is called.
3. Fires `chat.userMessage` to the extension.
4. Optimistically appends a user event to `session.events`.

**Receiving a Response:**
1. Extension sends `chat.event` messages.
2. `useChat` listener catches them.
3. Appends new events to `session.events`.
4. `Events` component re-renders.
5. `App.tsx` triggers auto-scroll.

## Key Directories

```text
src/
├── components/       # UI Components
│   ├── Events/       # Chat stream rendering (Messages, Tools)
│   ├── UserInput/    # Input editor and controls
│   ├── ui/           # Shadcn/UI primitives
│   ├── Header.tsx    # Top bar
│   └── ...
├── hooks/
│   └── useChat.tsx   # Main Context Provider & Logic
├── lib/
│   └── utils.ts      # Tailwind helpers (cn)
├── App.tsx           # Main Layout
├── main.tsx          # Entry point (Providers)
└── vscodeApi.ts      # VS Code API wrapper
```

## Common Tasks

- **Adding a new Tool UI:**
    - Modify `src/components/Events/ToolCall.tsx` to handle specific tool names or arguments visually.
- **Handling new Server Messages:**
    - Update `useEffect` listener in `src/hooks/useChat.tsx` to handle the new `message.type`.
    - Update state types if necessary.
- **Styling:**
    - Use Tailwind classes.
    - Check `src/index.css` for global styles.
    - Dark mode is default (VS Code theme awareness is partial/manual via colors like `bg-[#1e1e1e]`).

## Gotchas

- **File Extensions:** Use `.tsx` for components. The build will fail if JSX is in `.ts` files.
- **Imports:** Uses path alias `@/` for `src/`.
- **Event Coalescing:** The backend might send function calls and results as separate events. `components/Events/index.tsx` contains logic to visually merge them for a cleaner UI.
