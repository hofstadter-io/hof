# Planning: Filesystem Parity and Implementation

This plan outlines the steps to complete the filesystem implementation for the Veg agent, ensuring parity between the backend, agent tools, and the VS Code extension.

## 1. Backend Implementation (`lib/agent/runtime`)

- [x] **Environ Service Updates** (`lib/agent/runtime/services/environ`)
    - [x] Add `CreateDirectory` to `localEnviron` in `dagger_fs_mutate.go`.
    - [x] Ensure all mutation methods correctly increment tags and persist state.
- [x] **API Handler Implementation** (`lib/agent/runtime/api_environ.go`)
    - [x] Implement `fsWrite` (it's currently a stub).
    - [x] Implement `fsDelete` (it's currently a stub).
    - [x] Add `fsMkdir` handler.
    - [x] Add `fsRename` handler.
    - [x] Add `fsCopy` handler.
- [x] **Runtime Route Registration** (`lib/agent/runtime/runtime.go`)
    - [x] Register `POST /fs/mkdir`
    - [x] Register `POST /fs/rename`
    - [x] Register `POST /fs/copy`

## 2. Agent Tools Implementation (`lib/agent/tools/filesys`)

- [x] Add `FilesysMkdir` tool.
- [x] Add `FilesysRename` tool (Move).
- [x] Add `FilesysCopy` tool.
- [x] Ensure all mutation tools update `currEnv` in the agent state.

## 3. VS Code Extension Implementation (`extensions/vscode/extension/src/services`)

- [x] **`VegContentProvider` Implementation** (`filesystemProvider.ts`)
    - [x] Implement `writeFile`.
    - [x] Implement `createDirectory`.
    - [x] Implement `delete`.
    - [x] Implement `rename`.
    - [x] Implement `copy`.
- [x] **State Synchronization**
    - [x] Handle the response from mutation calls to update the environment URI/tag in the extension state.
    - [x] Trigger workspace folder updates or refreshes when the environment tag increments.

## 4. Verification

- [ ] Verify filesystem operations from the VS Code Explorer.
- [ ] Verify filesystem operations via Agent Tools.
- [ ] Verify tag incrementing and state persistence in the backend.
