# Plan - VS Code Virtual Filesystem Enhancements

## Goal
Improve the VS Code virtual filesystem (backed by Dagger) by adding "save back to dagger" functionality (uploading local files to the environment) and ensuring robust support for copy/paste operations. This involves refactoring filesystem operations into common handlers, updating the API to be session-aware, and handling global events for filesystem synchronization.

## Current Task(s)

- [ ] Add inline comments for the fs/scm uri's
- [ ] Ensure veg:// is only used in the vs code extension
- [ ] Ensure the implementation works as expected (with the human developer)

## Context
- **VS Code Virtual Filesystem**: Implemented in `extensions/vscode/extension/src/services/filesystemProvider.ts`.
- **Backend API**: Implemented in `lib/agent/runtime/handlers/api/filesys.go`.
- **Dagger Services**: Implemented in `lib/agent/services/environ/`.
- **Session State**: Session state (including `currEnv`) needs to stay in sync with filesystem modifications.
- **Validation**: The user must run the steps to validate changes

### Observed Trace (Tracer Bullet)
```
VegContentProvider.writeFile veg://host.docker.internal:5000/b5a1e081-9d9e-484c-96f0-9af13e3834f7%3A2/lib/agent/runtime/handlers/api/filesys.go {create: true, overwrite: true, unlock: false, atomic: false, append: false}
fsWrite: veg://host.docker.internal:5000/b5a1e081-9d9e-484c-96f0-9af13e3834f7%3A2/extensions/vscode/extension/src/services/filesystemProvider.ts /b5a1e081-9d9e-484c-96f0-9af13e3834f7:2/extensions/vscode/extension/src/services/filesystemProvider.ts
```

## Proposed Changes

### 1. Common Filesystem Handlers (`lib/agent/runtime/handlers/common/filesys.go`)
Create a new file to house filesystem logic that is shared between the REST API and potentially other entry points.
- Wrap `environ.Client()` calls.
- Accept an optional `sid` (session ID).
- **Corrected Path Logic**: Ensure the handler properly parses the URI and `path` query parameter to avoid redundant parsing and mismatched results.
- **Robustness**: Implement strict error handling for `SessionStatePut` and ensure `currEnv` is only updated on successful mutation.
- **Missing Handlers**: Implement `FilesysDiff` to centralize diffing logic.
- **Partial User Events**: Record a user event in the session history for all mutations. This event should capture the action (e.g., "manual file save") and the path, ensuring the session history reflects manual changes as well as agentic ones.
- Implement: `FilesysRead`, `FilesysWrite`, `FilesysList`, `FilesysStat`, `FilesysDelete`, `FilesysMkdir`, `FilesysCopy`, `FilesysRename`, `FilesysDiff`.

### 2. Refactor API Handlers (`lib/agent/runtime/handlers/api/filesys.go`)
- Update request structs to include `Sid`.
- Refactor all filesystem handlers to use the new `common` handlers.
- **Deduplication**: Remove redundant `url.Parse` and query parameter extraction logic that is now handled by `common`.
- Ensure `Runtime` is passed correctly.

### 3. VS Code Extension Improvements
- **Filesystem Provider (`filesystemProvider.ts`)**:
    - **Session Awareness**: Update `makeReq` calls to include `sid`.
    - **URI Scheme Resilience**: Ensure `updateUri` correctly handles `oci://` prefixes to prevent `veg://oci://` scheme corruption.
    - **Watch**: Implement the `watch` method by hooking into `extensionEmitter` events (e.g., `filesys.change` from backend).
    - **Upload Command**: Implement `veg.explorer.saveToVeg` (or `upload`) to allow saving local files/folders to the remote environment.
- **SCM Provider (`scmProvider.ts`)**:
    - **Merge to Veg**: Update `mergeDiff` to support merging changes *into* a `veg://` environment (target `veg://` URI), not just local files.

## Steps

### Step 1: Create Common Filesystem Handlers
- [x] Create `lib/agent/runtime/handlers/common/filesys.go`.
- [ ] Implement core functions (`Read`, `Write`, `List`, `Stat`, `Delete`, `Mkdir`, `Copy`, `Rename`, `Diff`).
    - [ ] Add missing `FilesysDiff` handler.
    - [ ] Implement robust error handling for `SessionStatePut`.
- [x] Add `SessionStatePut` logic to update `currEnv` on mutation.
- [x] Add session access verification using `SessionGet`.

### Step 2: Refactor API Handlers
- [x] Update `lib/agent/runtime/handlers/api/filesys.go` to use `common` handlers.
    - [x] Remove redundant path parsing (rely on `common` or `utils`).
- [x] Ensure `consts.VEG_USER_HEADER` is extracted and used.
- [x] Ensure `consts.VEG_DEFAULT_USER` is only used as a final fallback.

### Step 3: Enhance VS Code Extension (Filesystem)
- [x] Update `utils.ts` and `filesystemProvider.ts` to pass `sid` and `X-Veg-User` header.
    - [x] Update `websocket.ts` to pass `X-Veg-User` header during handshake.
    - [x] Fix `updateUri` to prevent `veg://oci://` scheme corruption.
- [x] Implement `veg.explorer.saveToVeg` command (upload logic).
- [x] Wire up `watch` using `extensionEmitter`.

### Step 4: Enhance VS Code Extension (SCM)
- [x] Update `mergeDiff` in `scmProvider.ts` to handle `veg://` destinations using `common.FilesysCopy` (or similar logic via API).

### Step 5: Debug & Fix Session Mismatches
- [x] **Analyze**: `extractSid` is identifying session IDs aggressively (e.g., from `library/alpine`), causing `SessionGet` to fail for valid OCI URIs.
- [x] **Verify**: Add logging to `extractSid` to confirm the hypothesis.
- [ ] **Fix**: 
    - [ ] Improve `extractSid` to differentiate between a session ID and a generic OCI path segment.
    - [ ] Resolve panic in `LookupEnviron` when handling non-session OCI images.
    - [ ] Ensure generic OCI reads work without error while maintaining session security for actual session URIs.

### Step 6: Verification

The user will build and verify as necessary. Ultimately this work will make that process more seamless.
