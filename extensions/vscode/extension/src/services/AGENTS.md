# Services Subsystem (`src/services`)

This directory contains the core services for the Veg extension: the **Virtual Filesystem Provider** and the **SCM (Source Control Management) Provider**. Together, they enable interacting with remote environments as if they were local workspace folders.

## Architecture: Push vs. Pull

The two providers follow fundamentally different interaction models:

- **SCM Provider (`scmProvider.ts`) - "Push"**: This provider is active and directive. When a diff is received from the backend, the extension explicitly creates `SourceControlResourceState` objects and "pushes" them into the VS Code SCM view. It has full control over the displayed tree.
- **Filesystem Provider (`filesystemProvider.ts`) - "Pull"**: This provider is reactive. VS Code's internal explorer driver "pulls" data by recursively calling `stat(root)` -> `readDirectory(root)` -> `stat(child)`. It relies on sequential, valid responses to construct the tree UI.

## URI and Path Handling

URI translation is handled in `utils.ts` by `vsUriToVeg`. 

### Key Constraints for Explorer Rendering:
1.  **Relative Rooting**: Because remote environments (like Dagger containers) use working directories, the `path` query parameter must use `./` as its base. The root of an environment is represented as `./`.
2.  **Stat Consistency**: For the explorer tree to render, `stat()` must succeed for every entry returned by `readDirectory()`. If a URI is malformed during this recursive pull, the tree will fail to render children even if the initial listing was successful.
3.  **Leaf Names**: `readDirectory` must return simple filenames (e.g., `main.go`), not absolute paths (e.g., `/app/main.go`), so that VS Code can correctly join them to the parent URI for the next sequential `stat` call.

## Filesystem Provider (`filesystemProvider.ts`)

The `VegContentProvider` is the sole and comprehensive implementation of `vscode.FileSystemProvider` for the `veg://` scheme. It manages command registration, URI translation, and communication with the extension backend API (`http://localhost:2257`).

### Core Data Structures and Types

```typescript
const SERVER_PORT = 2257;
const SERVER_URL = `http://localhost:${SERVER_PORT}`;

type Environ = {
	name?: string
	srcUri?: string
	srcPath?: string
	fromUri?: string
	dstPath?: string
	workdir?: string
}

type Folder = {
	uri: vscode.Uri
	sid: string
	name?: string | undefined
	base?: string | undefined
	session?: any
	environ?: Environ
}
type FolderListing = Array<[string, vscode.FileType]>
```

### URI Translation

The `vsUriToVeg` method is crucial for converting the VS Code internal URI representation (which uses path segments after the authority) into the backend's expected canonical URI format (which embeds the path in a query parameter).

```typescript
private vsUriToVeg(uri: vscode.Uri): vscode.Uri {
	// ... logic to parse 'veg://<auth>/<env>/<path>' into a base OCI-style URI
	// with a 'path' query parameter for the file path
}
```

### Command Registration

The file registers a large set of `veg.explorer.*` commands, primarily for file and environment management in the file explorer context menu:
- `veg.explorer.chat`: Initiates a new session focused on the selected file/directory.
- `veg.explorer.openEnviron`: Prompts the user to open a new directory/repo/image.
- `veg.explorer.showDiff`: Opens a side-by-side diff for all modified paths in the environment.
- `veg.explorer.mergeDiff`: Writes the modified, added, and deleted files from the virtual environment to a specified local `file://` target, utilizing `node:fs/promises`.

### Diff and Merge Implementation

The `showDiff` and `mergeDiff` commands fetch a unified diff payload from the backend (`/fs/diff`) for a given URI.

- `showDiff`: Iterates over `diff.modPaths` and opens multiple `vscode.diff` editors.
- `mergeDiff`: Uses the `diff` payload to perform native filesystem operations (`fs.writeFile`, `fs.rm`) on the destination URI (which must be `file://`).

```typescript
// mergeDiff logic snippet (uses node:fs/promises)
for (var path of diff.addPaths) {
  // ...
  const val = diff.files[path]
  const key = destination.path + path
  await fs.writeFile(key, val)
}
// ... similar for modPaths and delPaths
```
